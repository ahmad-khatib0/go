package server

// masterToProxyAsync forwards response from topic master to topic proxy
// in a fire-and-forget manner.
func (n *ClusterNode) masterToProxyAsync(msg *ClusterResp) error {
	// var unused bool
	// if c := n.callAsync("Cluster.TopicProxy", msg, &unused, nil); c.Error != nil {
	// 	return c.Error
	// }
	return nil
}
