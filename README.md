scgiclient
==========

I needed an scgi client so I could talk to rtorrent locally, without going through the hassle of setting up an http server in between. Especially since there's no Go lib that does xmlrpc which would allow me to write my own server for that.

This lib allows talking to both github.com/hoisie/web scgi server and rtorrent so it's good enough for me for now.

https://github.com/mpl/scgiclient/blob/master/example/rtorrent.go shows how to send a simple xml-rpc command to rtorrent.

