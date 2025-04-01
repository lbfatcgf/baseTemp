package tools

var stopCallbackList= make([]func(),0)

func AddOnStopSignal(callback func()){
	stopCallbackList = append(stopCallbackList, callback)
}

func StopSingalHandler(){
	for  _,callback := range stopCallbackList{
		callback()
	}
}