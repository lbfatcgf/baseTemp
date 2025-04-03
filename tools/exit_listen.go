package tools

var stopCallbackList= make([]func(),0)

// AddOnStopSignal 添加停止信号回调函数
func AddOnStopSignal(callback func()){
	stopCallbackList = append(stopCallbackList, callback)
}

// StopSingalHandler 停止信号处理函数
func StopSingalHandler(){
	for  _,callback := range stopCallbackList{
		callback()
	}
}