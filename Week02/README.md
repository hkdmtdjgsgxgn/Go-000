# Homework

* 问题： 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
* 答：需要Wrap这个错误，按照老师说的几个点，这属于使用其他库来协作，那么需要将根因用wrap保存。
* 建议：个别同学的作业我也看了下，我觉得就题论题而言，有些把各种web和sql实现全都放到作业代码里实际上是对作业本身场景缺乏更好的理解，这些场景应该是用最简单的方式构造一份出来更好，就像我这样☺️，好吧，我承认我自恋了。
# Slate

按照list记录下，我认为需要注意的点.

## Handling Error

* Assert errors for behaviour, not type
* 永远不要用panic去抛出业务逻辑错误，只有 越界、不可恢复的环境问题、栈溢出才用`panic`
* you only need to check the error value if you care about the result
* Eliminate error handling by eliminating errors
```go
type errWriter struct{
  io.Writer
  err error
}
func (e *errWriter)Write(buf []byte)(int, error){
  if e.err != nil{
    return 0, e.err
  }
  var n int
  n, e.err = e.Writer.Write(buf)
  return n, nil
}
func WriteResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
  ew := &errWriter{Writer: w}
  fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
  for _, h := range headers{
    fmt.Fprintf(ew, "%s: %s\r\n", h.Key, h.Value)
  }
  fmt.Fprint(ew, "\r\n")
  io.Copy(ew, body)
  return ew.err
}
```
* You should only handle errors once
* The error has been logged.
* The application is back to 100% integrity. 应用程序处理错误，保证100%完整性。
* The current error is not reported any longer.
* 直接返回错误，而不是每个错误产生的地方到处打日志
* 在程序顶部或者工作goroutine顶部（请求入口），适用`%+v`把堆栈详情记录
## Wrap errors

* wrap 适用场景：
    * 在应用层代码中使用`errors.New``errors.Errorf`
    * 调用其他包内函数直接返回。
    * 如果与其他库进行协作，应该用 wrap 把根因保存。
    * 自己写的基础库不应该 wrap error，只有你的应用程序代码选择wrap。
* **Packages that are reusable across many projects only return root error values**
    * 具有最高可用性的包只能返回根错误值。
* **If the error is not going to be handled, wrap and return up the call stack.**
    * 如果你不打算处理这个错误，那么用足够的上下文 wrap error 并返回到调用堆栈中
* **Once an error is handled, it is not allowed to be passed up the call stack any longer.**
    * 一旦确定函数或方法将处理错误，错误就不再是错误。如果函数或方法任然需要返回，则，它不能返回错误值，而应该返回 nil。
* 适用`errors.Cause`获取 root error, 再进行和`sentinel error`判定
# Go 1.13

## Customizing error tests with Is and As methods

![图片](https://uploader.shimo.im/f/ggG1zPrz2CoRZScZ.png!thumbnail?fileGuid=6RpJpp6kYqTxwyhx)

![图片](https://uploader.shimo.im/f/8ifvg3tkqeN0Kb3Z.png!thumbnail?fileGuid=6RpJpp6kYqTxwyhx)

# Reference

[https://dave.cheney.net/2012/01/18/why-go-gets-exceptions-right](https://dave.cheney.net/2012/01/18/why-go-gets-exceptions-right)

[https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux](https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux)

[https://dave.cheney.net/2014/11/04/error-handling-vs-exceptions-redux](https://dave.cheney.net/2014/11/04/error-handling-vs-exceptions-redux)

[https://rauljordan.com/2020/07/06/why-go-error-handling-is-awesome.html](https://rauljordan.com/2020/07/06/why-go-error-handling-is-awesome.html)

[https://morsmachine.dk/error-handling](https://morsmachine.dk/error-handling)

[https://blog.golang.org/error-handling-and-go](https://blog.golang.org/error-handling-and-go)

[https://www.ardanlabs.com/blog/2014/10/error-handling-in-go-part-i.html](https://www.ardanlabs.com/blog/2014/10/error-handling-in-go-part-i.html)

[https://www.ardanlabs.com/blog/2014/11/error-handling-in-go-part-ii.html](https://www.ardanlabs.com/blog/2014/11/error-handling-in-go-part-ii.html)

[https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

[https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)

[https://blog.golang.org/errors-are-values](https://blog.golang.org/errors-are-values)

[https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package](https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package)

[https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html](https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html)

[https://crawshaw.io/blog/xerrors](https://crawshaw.io/blog/xerrors)

[https://blog.golang.org/go1.13-errors](https://blog.golang.org/go1.13-errors)

[https://medium.com/gett-engineering/error-handling-in-go-53b8a7112d04](https://medium.com/gett-engineering/error-handling-in-go-53b8a7112d04)

[https://medium.com/gett-engineering/error-handling-in-go-1-13-5ee6d1e0a55c](https://medium.com/gett-engineering/error-handling-in-go-1-13-5ee6d1e0a55c)


