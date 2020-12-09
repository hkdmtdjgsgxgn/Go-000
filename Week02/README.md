# 学习笔记

## 作业

问题： 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
答：需要Wrap这个错误，按照老师说的几个点，这属于使用其他库来协作，那么需要将根因用wrap保存。

## 笔记

按照list记录下，我认为需要注意的点：

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

### Wrap errors

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
