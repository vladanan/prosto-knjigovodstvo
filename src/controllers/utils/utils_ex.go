package utils

// func ToStructHeaders(headers []byte) http.Header {
// 	var p http.Header
// 	err := json.Unmarshal(headers, &p)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return p
// }

// // func ToStructUser(user []byte) []models.User {
// // 	var p []models.User
// // 	err := json.Unmarshal(user, &p)
// // 	if err != nil {
// // 		log.Printf("json error: %v", err)
// // 	}
// // 	return p
// // }

// // bytearrayUser, e := json.Marshal(pgxUser)
// // if e != nil {
// // 	return l(r, e)
// // }

// func GetXForwardedFor(headers string) string {
// 	var headers_http http.Header = ToStructHeaders([]byte(headers))
// 	return headers_http["X-Forwarded-For"][0]
// }

// func HeadersToMap(headers []byte) map[string][]string {
// 	var h map[string][]string
// 	err := json.Unmarshal(headers, &h)
// 	if err != nil {
// 		// l(nil, 8, err)
// 		fmt.Println(err)
// 		return nil
// 	}
// 	return h
// }

// // fmt.Println("string concat rows:", pgxTests)
// // bytearray_tests, err2 := json.Marshal(pgx_tests)
// // if err2 != nil {
// // 	fmt.Printf("Json error: %v", err2)
// // }s
// // jsonstring_pitanja := string(bytearray_pitanja)
// // fmt.Println("json string pitanja:", jsonstring_pitanja)

// // https://stackoverflow.com/questions/54926712/is-there-a-way-to-list-keys-in-context-context
// func printContextInternals(ctx any, inner bool) {
// 	contextValues := reflect.ValueOf(ctx).Elem()
// 	contextKeys := reflect.TypeOf(ctx).Elem()

// 	if !inner {
// 		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
// 	}

// 	if contextKeys.Kind() == reflect.Struct {
// 		for i := 0; i < contextValues.NumField(); i++ {
// 			reflectValue := contextValues.Field(i)
// 			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

// 			reflectField := contextKeys.Field(i)

// 			if reflectField.Name == "Context" {
// 				printContextInternals(reflectValue.Interface(), true)
// 			} else {
// 				fmt.Printf("field name: %+v\n", reflectField.Name)
// 				fmt.Printf("value: %+v\n", reflectValue.Interface())
// 			}
// 		}
// 	} else {
// 		fmt.Printf("context is empty (int)\n")
// 	}
// }

// // al := clr.GetAuthLogger()
// // var Reset = "\033[0m"
// // log.SetFlags(log.LstdFlags | log.Lshortfile)
// // defer func() { log.SetFlags(log.LstdFlags); log.SetPrefix(Reset) }()

// // https://stackoverflow.com/questions/13765797/the-best-way-to-get-a-string-from-a-writer
// // ch := vezbamo.NewTestHandler(db.DB{})
// // rr := httptest.NewRecorder()
// // err := ch.HandleGetAll(rr, r)
// // err := vezbamo.GetTests(rr, r)
// // if err != nil {
// // 	// log.Println("greška na api")
// // 	templ.Handler(site.ServerError(clr.CheckErr(err))).Component.Render(context.Background(), w)
// // } else {
// // 	list_string := rr.Body.String() // r.Body is a *bytes.Buffer
// // 	dec := json.NewDecoder(strings.NewReader(list_string))
// // 	var all_tests []models.Test
// // 	if err := dec.Decode(&all_tests); err != nil {
// // 		// log.Println("greška json dekodera")
// // 		templ.Handler(site.ServerError(clr.CheckErr(err))).Component.Render(context.Background(), w)
// // 	} else {
// // 		templ.Handler(tests.List(all_tests)).Component.Render(context.Background(), w)
// // 	}
// // }

// // bytearray_headers, err2 := json.Marshal(r.Header)
// // if err2 != nil {
// // 	fmt.Printf("Sign_in: JSON error: %v", err2)
// // }
// // fmt.Print("\nSign_in: header:", string(bytearray_headers), "\n")
// // for item, index := range r.Header {
// // 	fmt.Print("\nSign_in: header:", item, index, "\n")
// // }
// // if already_authenticated, user_email, _, err := getSesionData(r); err == nil && already_authenticated {
// // 	_, data, _ := models.AuthenticateUser(user_email, "", already_authenticated, r)
// // 	if user, ok := data.(models.User); ok {
// // 		templ.Handler(dashboard.Dashboard(store, r, user)).Component.Render(context.Background(), w)
// // 	} else {
// // 		templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
// // 	}
// // } else if err != nil {
// // 	http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError) // greška u pristupu sesiji
// // 	return
// // } else {
// // 	templ.Handler(dashboard.Sign_in(store, r)).Component.Render(context.Background(), w)
// // }

// // 69 6E 74 65 72 6E 61 6C 20 75 73 65 72 20 65 72 72 6F 72
// // U+\x69\x6E\x74\x65\x72\x6E\x61\x6C\x20\x75\x73\x65\x72\x20\x65\x72\x72\x6F\x72
// // log.Println([]byte("internal user error"))
// // log.Println(string([]byte{105, 110, 116, 101, 114, 110, 97, 108, 32, 117, 115, 101, 114, 32, 101, 114, 114, 111, 114}))
// // log.Println("0x1 internal user error 0x02_9a U+0x69 0x6E 0x74 0x65 0x72 0x6E 0x61 0x6C 0x20 0x75 0x73 0x65 0x72 0x20 0x65 0x72 0x72 0x6F 0x72")
// // http.Error(w, "0o1 internal user error 0x02_9a U+0x69_6E_74_65_72_6E_61_6C_20_75_73_65_72_20_65_72_72_6F_72", http.StatusForbidden)

// // fmt.Println("\nr0:", r, r.Host, r.URL.Host, r.URL.User, r.URL.RequestURI())
// // fmt.Println("\nr1 hostr:", r.Host, "url host:", r.URL.Host, "url user:", r.URL.User, "url request uri:", r.URL.RequestURI())
// // fmt.Println("\nr2 header:", r.Header, "proto:", r.Proto, "remote addr:", r.RemoteAddr, "req uri:", r.RequestURI, "response:", r.Response, "url:", r.URL)
// // fmt.Println("accept lang:", r.Header.Get("Accept-Language"), r.URL.Path)

// var cl1 map[string]time.Time
// var bl1 map[string]time.Time

// func flodMiddleware1(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ip, port, err := net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			fmt.Printf("userip: %q is not IP:port", r.RemoteAddr)
// 		}
// 		fmt.Println("ip:", ip, "port:", port, "url:", r.URL, "path:", r.URL.Path, "xff:", r.Header.Get("X-Forwarded-For"), "rf:", r.Header.Get("Referer"))
// 		// 192.168.0.1
// 		// ruter server notbuk wifi 25, notbuk lan 226
// 		// telefon 21, stari notbuk 26
// 		// ruter 100.116.172.232
// 		// out 178.220.151.98
// 		// fmt.Println("\nr0:", r, r.Host, r.URL.Host, r.URL.User, r.URL.RequestURI())
// 		// fmt.Println("\nr1 hostr:", r.Host, "url host:", r.URL.Host, "url user:", r.URL.User, "url request uri:", r.URL.RequestURI())
// 		// fmt.Println("\nr2 header:", r.Header, "proto:", r.Proto, "remote addr:", r.RemoteAddr, "req uri:", r.RequestURI, "response:", r.Response, "url:", r.URL)
// 		fmt.Println("\nremote, cl, bl:", r.RemoteAddr, len(cl1), len(bl1))
// 		// fmt.Println("accept lang:", r.Header.Get("Accept-Language"), r.URL.Path)
// 		// fmt.Println("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"), r.URL.Path)
// 		enough := 0
// 		// log.Println("cl", len(cl))
// 		for _, c := range cl1 {
// 			// fmt.Println("<10", int(time.Since(c.ctime).Seconds()))
// 			if int(time.Since(c).Seconds()) < 7 {
// 				// log.Println("IMA u cl mlađih od n")
// 				enough++
// 			}
// 		}
// 		// fmt.Println("enough", enough)
// 		if enough == 0 {
// 			// log.Println("nema u cl mlađih od n i sve se briše")
// 			cl1 = make(map[string]time.Time)
// 			bl1 = make(map[string]time.Time)
// 		}
// 		nextcalls := 0
// 		if r.RemoteAddr != "" {
// 			remote := strings.Split(r.RemoteAddr, ":")[0]
// 			// log.Println("remote:", remote)
// 			// referer := strings.Split(r.Header.Get("Referer"), "//")[1]
// 			if _, found := bl1[remote]; found {
// 				fmt.Println("ima bl time > n", int(time.Since(bl1[remote]).Seconds()))
// 				if int(time.Since(bl1[remote]).Seconds()) > 3 {
// 					log.Println("prošlo više od n sec za bl del iz bl", len(bl1))
// 					delete(bl1, remote)
// 					next.ServeHTTP(w, r)
// 					nextcalls++
// 				} else {
// 					log.Println("malicious 0x1c0", remote)
// 					http.Error(w, "Internal user error 0x1c0", http.StatusForbidden)
// 				}
// 			} else if _, found := cl1[remote]; found {
// 				// log.Println("ima ga u cl")
// 				if float64(time.Since(cl1[remote]).Milliseconds()) < float64(1) {
// 					// log.Println("interval manji od 200ms ide u bl")
// 					bl1[remote] = time.Now()
// 				} else {
// 					// log.Println("interval veći od 200ms proverava se za cl > 1min", int(time.Since(cl[referer]).Minutes()))
// 					if int(time.Since(cl1[remote]).Seconds()) > 3 {
// 						// log.Println("interval za cl veći od n s sve okej i briše se iz cl")
// 						delete(cl1, remote)
// 						next.ServeHTTP(w, r)
// 						nextcalls++
// 					} else {
// 						// log.Println("interval za cl manji od 1m sve okej i updatuje se cl")
// 						cl1[remote] = time.Now()
// 						next.ServeHTTP(w, r)
// 						nextcalls++
// 					}
// 				}
// 			} else {
// 				// log.Println("nema ga u cl pa se dodaje")
// 				cl1[remote] = time.Now()
// 				next.ServeHTTP(w, r)
// 				nextcalls++
// 			}
// 		} else {
// 			// log.Println("nema ništa (nakon bl ili refresh), cl, bl", len(cl), len(bl))
// 			if len(bl1) > 0 {
// 				// log.Println("bl više od 0", bl)
// 				blstring := fmt.Sprint(bl1)
// 				blstring = strings.ReplaceAll(blstring, "map[", "")
// 				blstring = strings.ReplaceAll(blstring, "]", "")
// 				for _, t := range bl1 {
// 					blstring = strings.ReplaceAll(blstring, ":"+t.String(), "")
// 				}
// 				fmt.Println(blstring)
// 				blarray := strings.Split(blstring, " ")
// 				// kada ima jedan više od append jer počinje sa praznim pa se briše jedan
// 				// newarr = append(newarr[:0], newarr[1:]...)
// 				allokay := false
// 				var malicious string
// 				for _, a := range blarray {
// 					// fmt.Println("newarr", i, a)
// 					fmt.Println("time", int(time.Since(bl1[a]).Seconds()))
// 					if int(time.Since(bl1[a]).Seconds()) > 7 {
// 						log.Println("prošlo n s", a)
// 						delete(bl1, a)
// 						allokay = true
// 					} else {
// 						malicious = a
// 						allokay = false
// 					}
// 				}
// 				if allokay {
// 					next.ServeHTTP(w, r)
// 					nextcalls++
// 				} else {
// 					log.Println("malicious 0x2c0", malicious)
// 					http.Error(w, "Internal user error 0x2c0", http.StatusForbidden)
// 				}
// 			} else {
// 				next.ServeHTTP(w, r)
// 				nextcalls++
// 			}
// 		}
// 		// next.ServeHTTP(w, r)
// 		// nextcalls++
// 		// log.Println("cl:", cl)
// 		// log.Println("bl:", bl)
// 		// fmt.Println(r.URL.Path, nextcalls)
// 	})
// }

// // document.getElementById("kombi_gif").style.display = "none";
// // document.getElementById("kombi_muzika").style.display = "none";
// // document.getElementById("dugme_za_zadatak1").style.color = "white";
// // document.getElementById("dugme_za_zadatak1").style.display = "none";

// // let selectDOM = document.getElementById("lang");
// // let arr = [...selectDOM.children]
// // console.log("delay reloaddddddddd", arr[0])
// // arr.map(c => c.value == e.target.value ? c.setAttribute('selected', 'selected') : '');
