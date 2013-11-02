# Markdown Renderer

Markdown Renderer is a very simple HTTP server written in Go.  It
renders Markdown documents retrieved from another (specified) HTTP
server into HTML.

Markdown Renderer uses package
[`github.com/knieriem/markdown`](https://github.com/knieriem/markdown)
to render Markdown documents.  It can apply a CSS with the output
HTML.  You can write your own CSS, or download one from the Internet,
like [this](http://kevinburke.bitbucket.org/markdowncss/).

### A Use Case: Render Markdown Documents in SVN

This is a real use case, and the one that motivated me to write
Markdown Renderer.

In our company, we have an SVN server, on which our code and documents
reside.  We would like to be able to browse our documents from the Web
browser, in particular, we want those documents in Markdown syntax
being renderred to HTML.  However, the SVN server is not smart enough
to render Markdown documents; more than that, it does not even
recognizes file types of documents and returns all documents with
`Content-Type: text/plain` anyway.

This inspires me to set up an Nginx server, which `proxy_pass`es all
requests to the SVN server, and set the correct `Content-Type` by the
file extension name of corresponding document.  This can be done using
the `more_set_headers` directive provided by Nginx module
[`HttpHeadersMoreModule`](http://wiki.nginx.org/HttpHeadersMoreModule).
Any example Nginx configuration should be like this:

    server {
        location ~ \.docx$ {
            more_set_headers application/msword;
        }
        location ~ \.xlsx$ {
            more_set_headers application/vnd.ms-excel;
        }
    }

However, this module is not able to render Markdown text into HTML.
Indeed, I cannot find an Nginx module that can do this.  I tried to
write one by my own; however, had I digged into this work could I
realise what a pain it is to write an Nginx filter module!  This made
me to resort to an alternative way, to write a separate HTTP server,
instead of an Nginx module.  Thus comes Markdown Renderer.

With Markdown Renderer, a new `location` line can be added to above
example configuration:

        location ~ \.md$ {
            proxy_pass http://localhost:8002;
        }

where `localhost:8002` is supposed to be the Markdown Renderer server
started with proper command line flags set.  For example:

     ./markdown-renderer -addr=:8002 -data="http://svn-server:9006 -css="/markdown.css"

where `svn-server:9006` is just a replaceholder; you should change it
to your SVN or document server.

### Play with Markdown Renderer

The `nginx.conf` attached with this project configures two Nginx
virtual servers: the document-type-recognizer server as described in
above use case, and one that mimics the SVN/document server.

The recognizer server listens on `localhost:8001`, the Markdown
Renderer server listens on `localhost:8002`, and the fake SVN server
listens on `localhost:8003`.  They work in a chain:

     |browser|----|:8001|----(.md files)----|:8002|----|:8003|
                         \---(other docs)-------------/

If you want to setup this configuration on your computer and play with
it, these are the steps:

  1. Checkout and build Markdown Renderer:

        export ~/Projects/markdown-renderer
        cd ~/Projects
        go get github.com/wangkuiyi/markdown-renderer

  1. Download, build and install Nginx.

  1. Make Nginx use the configuration file provided with Markdown Renderer.

        cd /usr/local/nginx/conf  # suppose that Nginx was installed here.
        mv nginx.conf nginx.conf.bak  # backup the configuration file.
        ln -s ~/Projects/markdown-renderer/src/github.com/wangkuiyi/markdown-renderer/nginx.conf

  1. (Optional) Edit `nginx.conf` to specify the document root
     directory to be where Markdown Renderer source code is.

        location / {
            root   /Users/wangyi/Projects/markdown-renderer/src/github.com/wangkuiyi/markdown-renderer;
        }


  1. Start Nginx.

        /usr/local/nginx/sbin/nginx

  1. Build and start Markdown Renderer

        cd ~/Projects/markdown-renderer/src/github.com/wangkuiyi/markdown-renderer
        go install
        ~/Projects/markdown-renderer/bin/markdown-renderer


### Trouble Shooting

Markdown Renderer requires that the Markdown filename matches the
regular expression `^/([_a-zA-Z0-9]+)\\.md$)`.
