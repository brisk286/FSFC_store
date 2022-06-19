// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a git
// license that can be found in the LICENSE file.

package router

import (
	"fsfc_store/rpc/codec"
	"fsfc_store/rpc/compressor"
	"fsfc_store/rpc/serializer"
	"io"
	"net/rpc"
)

// Client rpc client based on net/rpc implementation
type Client struct {
	*rpc.Client
}

//Option provides options for rpc
type Option func(o *options)

type options struct {
	compressType compressor.CompressType
	serializer   serializer.Serializer
}

// WithCompress set client compression format
func WithCompress(c compressor.CompressType) Option {
	return func(o *options) {
		o.compressType = c
	}
}

// WithSerializer set client serializer
func WithSerializer(serializer serializer.Serializer) Option {
	return func(o *options) {
		o.serializer = serializer
	}
}

// NewClient Create a new rpc client
func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	options := options{
		compressType: compressor.Raw,
		serializer:   serializer.Proto,
	}
	for _, option := range opts {
		option(&options)
	}
	//rpc.NewClientWithCodec 创建使用自定义编解码器的client
	return &Client{rpc.NewClientWithCodec(
		codec.NewClientCodec(conn, options.compressType, options.serializer))}
}

// Call synchronously calls the rpc function
// 同步call
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.Client.Call(serviceMethod, args, reply)
}

// AsyncCall asynchronously calls the rpc function and returns a channel of *rpc.Call
func (c *Client) AsyncCall(serviceMethod string, args interface{}, reply interface{}) chan *rpc.Call {
	return c.Go(serviceMethod, args, reply, nil).Done
}
