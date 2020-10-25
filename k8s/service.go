package k8s

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// WriteLoop loop update map list
func WriteLoop(m map[string]string, mux *sync.RWMutex, ch <-chan Msg) {
	//log.Printf("event: %s\n", a)
	//log.Printf("map: %s\n", m)
	for {
		a := <-ch
		switch a.Type {
		case "ADDED":
			//log.Println("ADDED")
			for k, v := range a.Content {
				log.Printf("%s was added to the detection of the lists\n", k)
				mux.Lock()
				m[k] = v
				mux.Unlock()
			}
		case "DELETED":
			for k := range a.Content {
				log.Printf("%s was deleted from the detection of the lists\n", k)
				mux.Lock()
				delete(m, k)
				mux.Unlock()
			}
			log.Println(m)
		default:
			log.Println("The Service is being monitores....")
		}
	}
}

// ReadLoop loop http requst
func ReadLoop(m map[string]string, mux *sync.RWMutex, wg *sync.WaitGroup, exceptMsg chan ExceptMsg) {
	for {
		mux.RLock()
		for k, v := range m {
			url := v
			// log.Println(url)
			wg.Add(1)
			go func(k, url string) {
				defer wg.Done()
				var msg ExceptMsg
				var data string
				msg.Start = time.Now()
				resp, err := http.Get(m[k])
				//time.Sleep(3 * time.Second)
				msg.End = time.Now()
				// resolve resp nil to panic
				if resp != nil {
					defer resp.Body.Close()
				}
				// hanlder http get exception
				if err != nil {
					data = fmt.Sprintf("err: %s", err)
					msg.Message = fmt.Sprintf("\n\n > %s\n\n", data)
					msg.Service = k
					exceptMsg <- msg
					return
				}

				log.Println(k + " is ok...")

				// Body

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("read failed: %s", err)
				}
				// handler retrive result costs too long time
				if msg.End.Sub(msg.Start).Seconds() > 3 {
					msg.Message = fmt.Sprintf("%s\n\n", string(body))
					msg.Service = k
					exceptMsg <- msg
					return
				}

				defer resp.Body.Close()

			}(k, url)
		}
		wg.Wait()
		mux.RUnlock()
		time.Sleep(5 * time.Second)
	}
}

// LiveLoop loop get service lists
func LiveLoop(clientset *kubernetes.Clientset, ns, labels string, ch chan Msg) {
	var dst = Msg{}
	log.Println("Living loop Now")
	watcher, err := clientset.CoreV1().Services(ns).Watch(metav1.ListOptions{LabelSelector: labels})
	if err != nil {
		log.Printf("err: %s\n", err)
	}
	for {
		event := <-watcher.ResultChan()
		//fmt.Println(event)
		s := event.Object.(*v1.Service)
		//log.Println(event)
		switch event.Type {
		case "ADDED":
			dst.Type = "ADDED"
			data := make(map[string]string)
			data[s.Name] = "https://" + s.Name + ".domian.co/core/healthz.api"
			dst.Content = data
			ch <- dst
		case "DELETED":
			dst.Type = "DELETED"
			data := make(map[string]string)
			data[s.Name] = "https://" + s.Name + ".domain.co/core/healthz.api"
			dst.Content = data
			ch <- dst
		}
	}
}
