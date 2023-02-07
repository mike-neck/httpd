# httpq

http/html client with query selector.

# usage

"httpq" has an internal HTTP client.
The "-url" option specifies the URL to get the HTML.
You can query the elements using the query selector specified by the "-query" option
with the HTML obtained by the internal HTTP client as the input source.
For the obtained elements, specify the value of the desired attribute in "-values".

Example:

```bash
httpq -url https://www.google.com/ \
  -query 'a.gb1' \
  -values 'href'
```

The result of the above command would look something like this:

```text
https://www.google.co.jp/imghp?hl=ja&tab=wi
https://maps.google.co.jp/maps?hl=ja&tab=wl
https://play.google.com/?hl=ja&tab=w8
https://www.youtube.com/?tab=w1
https://news.google.com/?tab=wn
https://mail.google.com/mail/?tab=wm
https://drive.google.com/?tab=wo
https://www.google.co.jp/intl/ja/about/products?tab=wh
```


# Installation

T.B.D.
