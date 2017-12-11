package p2p

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"code.aliyun.com/chain33/chain33/common/crypto"
	"code.aliyun.com/chain33/chain33/queue"
	"code.aliyun.com/chain33/chain33/rpc"
	pb "code.aliyun.com/chain33/chain33/types"
	"golang.org/x/net/context"
	pr "google.golang.org/grpc/peer"
)

type p2pServer struct {
	imtx     sync.Mutex
	q        *queue.Queue
	book     *AddrBook
	nodeinfo *NodeBase
	InBound  map[string]*pb.P2PVersion
	OutBound map[string]*peer
}

func (s *p2pServer) addInBound(in *pb.P2PVersion) {
	defer s.imtx.Unlock()
	s.imtx.Lock()
	s.InBound[in.AddrFrom] = in
}

func (s *p2pServer) deleteInBound(addr string) {
	defer s.imtx.Unlock()
	s.imtx.Lock()
	if _, ok := s.InBound[addr]; ok {
		delete(s.InBound, addr)
	}

}

func (s *p2pServer) getBounds() []string {
	var inbounds = make([]string, 0)

	defer s.imtx.Unlock()
	s.imtx.Lock()
	for addr, _ := range s.InBound {
		inbounds = append(inbounds, addr)
	}

	return inbounds

}

func (s *p2pServer) Monitor() {
	go func() {
		for {
			for addr, peerinfo := range s.InBound {
				if (time.Now().Unix() - peerinfo.GetTimestamp()) > 120 {
					s.deleteInBound(addr)

				}
				time.Sleep(time.Second)
			}
		}
	}()
}

func (s *p2pServer) update(peer string) {
	defer s.imtx.Unlock()
	s.imtx.Lock()

	if peerinfo, ok := s.InBound[peer]; ok {
		peerinfo.Timestamp = time.Now().Unix()
	}

}

func NewP2pServer() *p2pServer {
	return &p2pServer{
		InBound: make(map[string]*pb.P2PVersion),
	}

}

func (s p2pServer) checkSign(in *pb.P2PPing) bool {
	data := pb.Encode(in)
	sign := in.GetSign()
	if sign == nil {
		return false
	}
	c, err := crypto.New(pb.GetSignatureTypeName(int(sign.Ty)))
	if err != nil {
		return false
	}
	pub, err := c.PubKeyFromBytes(sign.Pubkey)
	if err != nil {
		return false
	}
	signbytes, err := c.SignatureFromBytes(sign.Signature)
	if err != nil {
		return false
	}
	return pub.VerifyBytes(data, signbytes)
}

func (s p2pServer) Ping(ctx context.Context, in *pb.P2PPing) (*pb.P2PPong, error) {
	log.Debug("PING", "RECV PING", in)
	getctx, ok := pr.FromContext(ctx)
	if ok {
		log.Debug("PING addr", "Addr", getctx.Addr.String())
		remoteaddr := fmt.Sprintf("%s:%v", in.Addr, in.Port)

		log.Debug("RemoteAddr", "Addr", remoteaddr)
		if s.checkSign(in) == true {
			//TODO	s.update(remoteaddr)
			log.Debug("PING CHECK SIGN SUCCESS")
		}

	}

	s.update(fmt.Sprintf("%v:%v", in.Addr, in.Port))
	log.Debug("Send Pong", "Nonce", in.GetNonce())
	return &pb.P2PPong{Nonce: in.GetNonce()}, nil

}

// 获取地址
func (s *p2pServer) GetAddr(ctx context.Context, in *pb.P2PGetAddr) (*pb.P2PAddr, error) {
	log.Debug("GETADDR", "RECV ADDR", in, "OutBound Len", len(s.OutBound))
	var addrlist = make([]string, 0)
	if len(s.OutBound) != 0 {
		for peeraddr, _ := range s.OutBound {
			addrlist = append(addrlist, peeraddr)
		}
	}

	for peeraddr, _ := range s.book.addrPeer {
		addrlist = append(addrlist, peeraddr)
		if len(addrlist) > MaxAddrListNum { //最多一次性返回256个地址
			break
		}
	}

	return &pb.P2PAddr{Nonce: in.Nonce, Addrlist: addrlist}, nil
}

// 版本
func (s *p2pServer) Version(ctx context.Context, in *pb.P2PVersion) (*pb.P2PVerAck, error) {

	remoteaddr := in.AddrRecv
	//localNetwork, _ := NewNetAddressString(localaddr)
	remoteNetwork, _ := NewNetAddressString(remoteaddr)

	log.Debug("RECV PEER VERSION", "VERSION", *in)
	s.book.addAddress(remoteNetwork)
	return &pb.P2PVerAck{Version: Version, Service: 6, Nonce: in.Nonce}, nil
}
func (s *p2pServer) Version2(ctx context.Context, in *pb.P2PVersion) (*pb.P2PVersion, error) {

	getctx, ok := pr.FromContext(ctx)
	var peeraddr string
	if ok {
		peeraddr = strings.Split(getctx.Addr.String(), ":")[0]
		log.Debug("Peer addr", "Addr", peeraddr)

	}
	//in.AddrFrom 表示远程客户端的地址,如果客户端的远程地址与自己定义的addrfrom 地址一直，则认为在外网
	if strings.Split(in.AddrFrom, ":")[0] == peeraddr {
		remoteNetwork, err := NewNetAddressString(in.AddrFrom)
		if err == nil && in.GetService() == NODE_NETWORK+NODE_GETUTXO+NODE_BLOOM {
			s.book.addAddress(remoteNetwork)
		}
	}

	log.Debug("RECV PEER VERSION", "VERSION", *in)
	if in.Version > Version {
		//Not support this Version
		log.Error("VersionCheck", "Error", "Version not Support")
		return nil, fmt.Errorf("Version No Support")
	}
	s.addInBound(in)
	//addrFrom:表示自己的外网地址，addrRecv:表示对方的外网地址
	return &pb.P2PVersion{Version: Version, Service: SERVICE, Nonce: in.Nonce,
		AddrFrom: in.AddrRecv, AddrRecv: fmt.Sprintf("%v:%v", peeraddr, strings.Split(in.AddrFrom, ":")[1])}, nil
}

//grpc 接收广播交易
func (s *p2pServer) BroadCastTx(ctx context.Context, in *pb.P2PTx) (*pb.Reply, error) {
	log.Debug("RECV TRANSACTION", "in", in)
	//发送给消息队列Queue
	client := s.q.GetClient()
	msg := client.NewMessage("mempool", pb.EventTx, in.Tx)
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	return resp.GetData().(*pb.Reply), nil
}

func (s *p2pServer) GetBlocks(ctx context.Context, in *pb.P2PGetBlocks) (*pb.P2PInv, error) {
	if in.GetEndHeight()-in.GetStartHeight() > 100 {
		return nil, errors.New("out of range")
	}

	//TODO GetHeaders
	client := s.q.GetClient()
	msg := client.NewMessage("blockchain", pb.EventGetHeaders, &pb.ReqBlocks{Start: in.StartHeight, End: in.EndHeight,
		Isdetail: false})
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}

	if resp.Err() != nil {
		return nil, resp.Err()
	}

	headers := resp.Data.(*pb.Headers)
	var invs = make([]*pb.Inventory, 0)
	for _, item := range headers.Items {
		var inv pb.Inventory
		inv.Ty = MSG_BLOCK
		//inv.Hash = item.Block.Hash()
		inv.Height = item.GetHeight()
		invs = append(invs, &inv)
	}
	return &pb.P2PInv{Invs: invs}, nil
}

//服务端查询本地mempool
func (s *p2pServer) GetMemPool(ctx context.Context, in *pb.P2PGetMempool) (*pb.P2PInv, error) {
	log.Debug("GetMempool", "version", in)
	client := s.q.GetClient()
	msg := client.NewMessage("mempool", pb.EventGetMempool, nil)
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	var invlist = make([]*pb.Inventory, 0)
	txlist := resp.GetData().(*pb.ReplyTxList)
	for _, tx := range txlist.GetTxs() {

		invlist = append(invlist, &pb.Inventory{Hash: tx.Hash(), Ty: MSG_TX})
	}

	return &pb.P2PInv{Invs: invlist}, nil
}

func (s *p2pServer) GetData(ctx context.Context, in *pb.P2PGetData) (*pb.InvDatas, error) {
	log.Debug("GetDataTx", "p2p version", in.GetVersion())
	var p2pInvData = make([]*pb.InvData, 0)
	//先获取本地mempool 模块的交易
	client := s.q.GetClient()
	msg := client.NewMessage("mempool", pb.EventGetMempool, nil)
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	txlist := resp.GetData().(*pb.ReplyTxList)
	txs := txlist.GetTxs()
	var txmap = make(map[string]*pb.Transaction)
	for _, tx := range txs {
		txmap[hex.EncodeToString(tx.Hash())] = tx
	}

	invs := in.GetInvs()
	for _, inv := range invs { //过滤掉不需要的数据
		var invdata pb.InvData
		if inv.GetTy() == MSG_TX {
			txhash := hex.EncodeToString(inv.GetHash())
			if tx, ok := txmap[txhash]; ok {

				invdata.Value = &pb.InvData_Tx{Tx: tx}
				invdata.Ty = MSG_TX
				p2pInvData = append(p2pInvData, &invdata)
			}

		} else if inv.GetTy() == MSG_BLOCK {
			cli := rpc.NewClient("channel", "")
			cli.SetQueue(s.q)
			height := inv.GetHeight() //TODO
			blocks, err := cli.GetBlocks(height, height, false)
			if err != nil {
				log.Error("GetBlocks Err", "Err", err.Error())
				continue
			}
			for _, item := range blocks.Items {
				invdata.Ty = MSG_BLOCK
				invdata.Value = &pb.InvData_Block{Block: item.Block}
				p2pInvData = append(p2pInvData, &invdata)
			}

		}

	}

	return &pb.InvDatas{Items: p2pInvData}, nil
}

func (s *p2pServer) GetHeaders(ctx context.Context, in *pb.P2PGetHeaders) (*pb.P2PHeaders, error) {
	log.Debug("GetHeaders", "p2p version", in.GetVersion())
	if in.GetEndHeigh()-in.GetStartHeight() > 2000 || in.GetEndHeigh() < in.GetStartHeight() {
		return nil, fmt.Errorf("out of range")
	}

	client := s.q.GetClient()
	msg := client.NewMessage("blockchain", pb.EventGetHeaders, pb.ReqBlocks{Start: in.GetStartHeight(), End: in.GetEndHeigh()})
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	headers := resp.GetData().(*pb.Headers)

	return &pb.P2PHeaders{Headers: headers.GetItems()}, nil
}

func (s *p2pServer) GetPeerInfo(ctx context.Context, in *pb.P2PGetPeerInfo) (*pb.P2PPeerInfo, error) {
	//log.Debug("GetPeerInfo", "p2p version", in.version)

	client := s.q.GetClient()
	msg := client.NewMessage("mempool", pb.EventGetMempoolSize, nil)
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	meminfo := resp.GetData().(*pb.MempoolSize)
	var peerinfo pb.P2PPeerInfo

	pub, err := Pubkey(s.book.key)
	if err != nil {
		log.Error("getpubkey", "error", err.Error())
	}

	//get header
	msg = client.NewMessage("blockchain", pb.EventGetLastHeader, nil)
	client.Send(msg, true)
	resp, err = client.Wait(msg)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	header := resp.GetData().(*pb.Header)

	peerinfo.Header = header
	peerinfo.Name = pub
	peerinfo.MempoolSize = int32(meminfo.GetSize())
	peerinfo.Addr = EXTERNALADDR
	peerinfo.Port = int32(s.nodeinfo.externalAddr.Port)
	return &peerinfo, nil
}

func (s *p2pServer) BroadCastBlock(ctx context.Context, in *pb.P2PBlock) (*pb.Reply, error) {
	client := s.q.GetClient()
	msg := client.NewMessage("blockchain", pb.EventBlockBroadcast, in.GetBlock())
	client.Send(msg, true)
	resp, err := client.Wait(msg)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	return resp.GetData().(*pb.Reply), nil
}