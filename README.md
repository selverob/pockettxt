PocketTxt
=========

Export articles in your [Pocket](http://getpocket.com) as text files

##What does it do?
PocketTxt is a utility that enables you to download your unread Pocket articles, saving them into a single text file that can be then uploaded to your Kindle or another ebook reader.

##How does it work?
It uses Pocket's API to authenticate and get URLs of your unread articles. Since Pocket's API to access scraped articles is not publicly available, those URLs are then passed to [Diffbot](http://diffbot.com).

##How do I use it?
Just `go get github.com/sellweek/pockettxt`, run the program and follow the on-screen instructions.

If you want a little bit more control, though, you can use some command-line flags.

* `-aToken=<token>` bypasses the authentication process and just uses the provided token. Your access token is printed while running pockettxt. It looks something like this: `0123456-789a-bcde-f123-456789`

*`-filename=<filename>` sets the name of the file which will be written
