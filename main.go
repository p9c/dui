package main

import (
	"encoding/json"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
	"log"
	"math/rand"
	"os"
	"time"
)

var html = `<!DOCTYPE html><html lang="en">
<head><meta charset="utf-8"><title>DUO</title>
<style>
html,body{
display:flex;
	width:100%;
	height:100vh;
	margin:0;
	padding:0;
	justify-content:center;
	align-items:center;
	text-align:center;
	background:#303030;
	color:#cfcfcf;
}</style></head><body><button id="kernel">DUO</button></body></html>`

var js = `var test = document.createElement("p");
var node = document.createTextNode("This is new.");
para.appendChild(node);
var element = document.getElementById("div1");
element.appendChild(para);`

var (
	qmlObjects = make(map[string]*core.QObject)

	qmlBridge          *QmlBridge
	manipulatedFromQml *widgets.QWidget

	colors = []string{"red", "green", "blue", "cyan", "magenta", "yellow", "gray"}
)

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	//app.Main()

	widgets.NewQApplication(len(os.Args), os.Args)
	w := widgets.NewQMainWindow(nil, 0)
	v := webengine.NewQWebEngineView(nil)
	v.SetHtml(html, core.NewQUrl())

	w.SetCentralWidget(v)

	w.Show()
	v.Page().RunJavaScript(js)
	widgets.QApplication_Exec()

}

func newCppWidget() *widgets.QWidget {

	var button = widgets.NewQPushButton2("Call Qml Function", nil)
	button.ConnectClicked(func(_ bool) {
		rand.Seed(time.Now().UnixNano())
		//qmlBridge.SendToQml("GoButton", "click", colors[rand.Intn(len(colors))])
	})

	manipulatedFromQml = widgets.NewQWidget(nil, 0)

	var layout = widgets.NewQVBoxLayout()
	layout.AddWidget(button, 0, 0)
	layout.AddWidget(manipulatedFromQml, 0, 0)

	var widget = widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	return widget
}

func newQmlWidget() *quick.QQuickWidget {
	var quickWidget = quick.NewQQuickWidget(nil)
	quickWidget.SetResizeMode(quick.QQuickWidget__SizeRootObjectToView)

	initQmlContext(quickWidget)

	quickWidget.SetSource(core.NewQUrl3("qrc:/qml/bridge.qml", 0))

	return quickWidget
}

func initQmlContext(quickWidget *quick.QQuickWidget) {

	var m = map[string]map[string]string{
		"QmlButton": {
			"color":        "lightGray",
			"pressedColor": "darkGray",
			"text":         "Call Go Function",
		},
	}

	var b, err = json.Marshal(m)
	if err != nil {
		log.Println("initQmlContext", err)
	}
	quickWidget.RootContext().SetContextProperty2("qmlInitContext", core.NewQVariant1(string(b)))
}

type QmlBridge struct {
	core.QObject

	_ func(source, action, data string) `signal:"sendToQml"`
	_ func(source, action, data string) `slot:"sendToGo"`

	_ func(object *core.QObject) `slot:"registerToGo"`
	_ func(objectName string)    `slot:"deregisterToGo"`
}
