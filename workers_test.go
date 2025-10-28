package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
	"github.com/joho/godotenv"
)

func TestWorker(t *testing.T) {
	godotenv.Load()
	StartWorkerPool(&ApiConfig{})

	cases := []struct {
		id         string
		resultChan chan YtDlpTaskResult
	}{
		{id: "diIFhc_Kzng", resultChan: make(chan YtDlpTaskResult)}, {id: "AE005nZeF-A", resultChan: make(chan YtDlpTaskResult)},
		{id: "uzS3WG6__G4", resultChan: make(chan YtDlpTaskResult)}, {id: "X_SEwgDl02E", resultChan: make(chan YtDlpTaskResult)},
		{id: "r4l9bFqgMaQ", resultChan: make(chan YtDlpTaskResult)}, {id: "HWDaIRe8_XI", resultChan: make(chan YtDlpTaskResult)},
		{id: "ncqkC9Ob2ZI", resultChan: make(chan YtDlpTaskResult)}, {id: "Dlz_XHeUUis", resultChan: make(chan YtDlpTaskResult)},
		{id: "P18g4rKns6Q", resultChan: make(chan YtDlpTaskResult)},
	}

	for _, c := range cases {
		mutex.Lock()
		pushTask(&Tasks, YtDlpTask{
			YoutubeId:  c.id,
			Priority:   0,
			ResultChan: c.resultChan,
		})
		// fmt.Println(Tasks)
		mutex.Unlock()
		go func() {
			workerResult := <-c.resultChan
			fmt.Println("Extracted Video :", workerResult.result.Title, "\t(good job worker)")
			fmt.Println("Youtube audio streaming url: ", youtube.GetAudioStreamingUrl(workerResult.result))
			fmt.Println(Tasks)
		}()
	}
	// Addition of a 'priority task'
	mutex.Lock()
	pushTask(&Tasks, YtDlpTask{
		YoutubeId:  "2npegbvmfso",
		Priority:   1,
		ResultChan: make(chan YtDlpTaskResult),
	})

	mutex.Unlock()
	mutex.Lock()
	pushTask(&Tasks, YtDlpTask{
		YoutubeId:  "DloZ1xZHCmo",
		Priority:   1,
		ResultChan: make(chan YtDlpTaskResult),
	})

	mutex.Unlock()

	time.Sleep(1 * time.Minute)
}

func TestYtUrlChecker(t *testing.T) {
	cases := []string{
		`https://rr3---sn-h50gpup0nuxaxjvh-hg06.googlevideo.com/videoplayback?expire=1761688513&ei=YOcAafLFOdTGi9oPi5TYyAo&ip=102.115.47.147&id=o-AG-1lt9LaS4-WmSprgWtumFfp9Mj0C9rJZqYL1DHpZWo&itag=140&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&cps=124&met=1761666912%2C&mh=gi&mm=31%2C29&mn=sn-h50gpup0nuxaxjvh-hg06%2Csn-woc7knez&ms=au%2Crdu&mv=m&mvi=3&pl=22&rms=au%2Cau&initcwndbps=1795000&bui=AdEuB5TtY9i-2FwaGKeHBSKA1RfJPLKe8s1OKSfSgjPWBSGzxZstAtzPXTLfC074vm8yorqdwS7NweyZ&spc=6b0G_G-KfKL_&vprv=1&svpuc=1&mime=audio%2Fmp4&rqh=1&gir=yes&clen=5084781&dur=314.119&lmt=1734769562910023&mt=1761666372&fvip=1&keepalive=yes&fexp=51552689%2C51565116%2C51565681%2C51580968&c=ANDROID&txp=4532434&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&sig=AJfQdSswRQIhAPWDyV5rXeJ1JxwujcCifk_WSeH_IqmOnxEBOFXTkR5nAiB40Amh2QQI25F3PsNXe65xCZvORH6iSNfphMC3CO0qBQ%3D%3D&lsparams=cps%2Cmet%2Cmh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Crms%2Cinitcwndbps&lsig=APaTxxMwRAIgEk3zyp_X3DUzDWVHOwCWQU--AmDgg4JrYMI3p_cZlHoCIDummepC8v5pAgaoSafCn4o66ocHwsLOCJTi-R83Lcae`,
		`https://rr1---sn-h50gpup0nuxaxjvh-hg0k.googlevideo.com/videoplayback?expire=176
0881815&ei=N5j0aMzDDNbCmLAPn4O9yQs&ip=102.115.32.145&id=o-ALGPGia3kyBqP5W9aooyCi
xhmGQylqZgUR6_w78epzSK&itag=140&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%
3D&met=1760860215%2C&mh=9y&mm=31%2C29&mn=sn-h50gpup0nuxaxjvh-hg0k%2Csn-woc7kn7y&
ms=au%2Crdu&mv=m&mvi=1&pl=20&rms=au%2Cau&initcwndbps=2061250&bui=ATw7iSWwQT_hwLy
VxoRGDpMnmarxSYahy_B6XIlT3cYbjhqAGDPVr-b7Vh-UgYYFpcCnnARPPaLyZb2h&vprv=1&svpuc=1
&mime=audio%2Fmp4&ns=4AsZh6R13LpEvCA6Ms6o56AQ&rqh=1&gir=yes&clen=3416026&dur=211
.022&lmt=1755677449145957&mt=1760859648&fvip=2&keepalive=yes&lmw=1&fexp=51557447
%2C51565115%2C51565682%2C51580970&c=TVHTML5&sefc=1&txp=6208224&n=XIcdDsfftfn37g&
sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cvprv%2C
svpuc%2Cmime%2Cns%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&lsparams=met%2Cmh%2Cmm%2Cmn%2Cm
s%2Cmv%2Cmvi%2Cpl%2Crms%2Cinitcwndbps&lsig=APaTxxMwRAIgWAudXtODVNLsYkWEIpdwcJ2tg
gfA0iLLaR0PMo6Un9MCICmXa_By71w4NbVmNWEy4uXQ-0tntsrvJNaS-14qUG8W&sig=AJfQdSswRQIhAI9HSsTLSwBrgB00j1M0kT4y5CoM4YDNZbVQMvNogGXWAiAbUzzXs7g07ij6I-4SN5AG99C1p9jxNaekGTx1N3fxOQ%3D%3D`,
		`https://rr2---sn-h50gpup0nuxaxjvh-hg0d.googlevideo.com/videoplayback?expire=1760881816&ei=OJj0aMv1Fo6vp-oPvsSjuQM&ip=102.115.32.145&id=o-AD4tYj3C5zGWLVDHWkV0us8XKSzUrDb-LBiABCHI853i&itag=140&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&met=1760860216%2C&mh=YK&mm=31%2C29&mn=sn-h50gpup0nuxaxjvh-hg0d%2Csn-4g5ednks&ms=au%2Crdu&mv=m&mvi=2&pl=20&rms=au%2Cau&initcwndbps=2153750&bui=ATw7iSXThzRLi7zf3_ADUXmRPzriR-F9RXYClqubeAF-qRtTs0p4U8Lu_YxBUdY8gCKl0rCWF3i7tOcm&vprv=1&svpuc=1&mime=audio%2Fmp4&ns=frwQ1YAQYIy4_u-xgSc5Ks0Q&rqh=1&gir=yes&clen=4179169&dur=258.182&lmt=1760371856312450&mt=1760859889&fvip=4&keepalive=yes&lmw=1&fexp=51557447%2C51565116%2C51565681%2C51580970&c=TVHTML5&sefc=1&txp=6208224&n=iqoZDDDM9Puz5Q&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cvprv%2Csvpuc%2Cmime%2Cns%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&lsparams=met%2Cmh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Crms%2Cinitcwndbps&lsig=APaTxxMwRQIhANIZvUGVfc2N9nFT0-_irp68Gr5XmfAqnzjC0PYIRmGUAiAsyX2tUGB4GiiR_gUEzcj1IRHJXSCdwobt9KuFoUk8Mg%3D%3D&sig=AJfQdSswRQIhAJW8OAjugsvSdDJqRoK8xW-GQIJctkY-3SCKgVI2gXG4AiB6qJqqOtvuB0q808h7JW-3CT5Zl_GdmHzl3R5j_fsYqg%3D%3D`,
		`https://rr4---sn-h50gpup0nuxaxjvh-hg0k.googlevideo.com/videoplayback?expire=1761438888&ei=SBj9aJPcEM7V0u8P8crYuAw&ip=102.115.16.96&id=o-ANsVa2NEv45VoyCXyHHmIZaufM9EhvaXwDsPAa3xSrxH&itag=140&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&met=1761417288%2C&mh=xQ&mm=31%2C29&mn=sn-h50gpup0nuxaxjvh-hg0k%2Csn-woc7knel&ms=au%2Crdu&mv=m&mvi=4&pl=23&rms=au%2Cau&initcwndbps=1652500&bui=AdEuB5R-ZWlkUCjlNG5dtYqZ8bh8KJi2Dvmig5WlwaGd0MH48trmd1SZ_mVZwYeNSqAWukHQDxuy-lFc&spc=6b0G_K4HswAt&vprv=1&svpuc=1&mime=audio%2Fmp4&rqh=1&gir=yes&clen=2803746&dur=173.174&lmt=1761388704127772&mt=1761416748&fvip=2&keepalive=yes&fexp=51552689%2C51565115%2C51565682%2C51580968&c=ANDROID&txp=4532534&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cbui%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&sig=AJfQdSswRAIgCccXIn816FSdpB68jYTObHaHx69yx6huFwG92A5GK4gCIAn0gPFJ8PZJgxIX7ucFjP03mJMkWf6kqE5wi1wKGcUy&lsparams=met%2Cmh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Crms%2Cinitcwndbps&lsig=APaTxxMwRgIhALgJkED1hWgG4yUGokO7ivlpvl-BJrJX0xKuEchWrAl_AiEA75FcDnp4mxqC9cp-kWYw3NFeF8PF2JpcMrBfGtH3c-8%3D`,
	}
	for _, c := range cases {
		checkYoutubeUrlResponse(c)
	}
}

// func TestDownload(t *testing.T) {

// 	_, err := downloadFile("hehe")

// 	if err != nil {
// 		t.Error(err)
// 	}

// }
