# LINO

The Line notification library for Golang.

- Line Notify API Document
    - [English](https://notify-bot.line.me/doc/en/)
    - [日本語](https://notify-bot.line.me/doc/ja/)

## Install

```sh
$ go get -u github.com/tishibas/lino
```

## Support

Currently, this library only supports Notify(POST `https://notify-api.line.me/api/notify`).



## Hello, world

Here is a simple "Hello, world" example for LINE notification.

```go
c := lino.New(&Config{
		AccessToken: "<ACCESS_TOKEN>",
    })

// only message
c.Notify(&lino.RequestNotify{
    Message: "message",
})

// message with images
imageThumbnail := "https://example.com/foo.jpg"
imageFullsize := "https://exmaple.com/bar.jpg"
c.Notify(&lino.RequestNotify{
    Message:        "message",
    ImageThumbnail: &imageThumbnail
    ImageFullsize:  &imageFullsize,
})

```

