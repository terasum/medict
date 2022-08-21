package main

import (
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/internal/static"
	"github.com/terasum/medict/pkg/apis"
	"github.com/webview/webview"
)

const html = `<button id="increment">Tap me</button>
<div>You tapped <span id="count">0</span> time(s).</div>
<script>
  const [incrementElement, countElement] =
    document.querySelectorAll("#increment, #count");
  document.addEventListener("DOMContentLoaded", () => {
    incrementElement.addEventListener("click", () => {
      window.increment().then(result => {
        countElement.textContent = result.count;
      });
    });
  });
</script>`

type IncrementResult struct {
	Count uint `json:"count"`
}

func main() {
	var count uint = 0
	w := webview.New(false)

	config, err := entry.DefaultConfig()
	if err != nil {
		panic(err)
	}
	dicts, err := apis.NewDictsAPI(config)
	if err != nil {
		panic(err)
	}
	defer dicts.Destroy()
	//staticInfos := apis.NewStaticInfosApi(config)

	go static.StartStaticServer(config)

	defer w.Destroy()

	w.SetTitle("Medict App")
	w.SetSize(800, 600, webview.HintNone)

	w.Bind("increment", func() IncrementResult {
		count++
		return IncrementResult{Count: count}
	})
	w.Navigate("http://localhost:19191/")
	w.Run()
}
