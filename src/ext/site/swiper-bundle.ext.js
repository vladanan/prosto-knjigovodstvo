function pageReload() {
  location.reload()
}

function delayReload() {
  setTimeout(pageReload(), 200);
}

function sendLang(e) {
  const xhttp = new XMLHttpRequest();
  xhttp.onload = function () {
    // console.log("function name:", e.target.value)
    // document.getElementById("ElementById").innerHTML = xhttp.responseText;
  }
  // za mozila fajerfoks, mora da stavi asinc ovde da bude false i onda se žali ali radi, ostali browseri rade u oba slučaja bez ikakvih primedbi
  let url = ""
  // test your target here
  // console.log("function name:", e.target.value)
  e.target.value == "sign_in" ? url = "" : url = e.target.value
  xhttp.open("POST", "/" + url, false);
  xhttp.send();
  delayReload()
}

function confirmData(e) {
  document.getElementById("fk_KjUkYoWMqzciQIhk5i4lYeqS2ewYAbpaBae3078ytRTtxvm8h_fJi").setAttribute("value", "");
  document.getElementById("fk_KjUkYoWMqzciQIhk5i4lYeqS2ewYAbpaBae3078ytRTtxvm8h_fJi").removeAttribute("required");
}

function showHidePassword0(e) {
  span = document.getElementById("show-hide-pass0")
  // attrib = span.getAttribute("src")
  state = span.getAttribute("state")
  // console.log()
  if (state == "hide") {
    // span.setAttribute("src", "static/site/show-password.png")
    span.innerHTML = "&#9678;"
    span.setAttribute("state", "show")
    document.getElementById("password0").setAttribute("type", "text")
  } else {
    // span.setAttribute("src", "static/site/hide-password.png")
    span.innerHTML = "&#9673;"
    span.setAttribute("state", "hide")
    document.getElementById("password0").setAttribute("type", "password")
  }
}

function showHidePassword(e) {
  span = document.getElementById("show-hide-pass")
  // attrib = span.getAttribute("src")
  state = span.getAttribute("state")
  // console.log()
  if (state == "hide") {
    // span.setAttribute("src", "static/site/show-password.png")
    span.innerHTML = "&#9678;"
    span.setAttribute("state", "show")
    document.getElementById("password1").setAttribute("type", "text")
    if (document.getElementById("password2") != null) {
      document.getElementById("password2").setAttribute("type", "text")
    }
  } else {
    // span.setAttribute("src", "static/site/hide-password.png")
    span.innerHTML = "&#9673;"
    span.setAttribute("state", "hide")
    document.getElementById("password1").setAttribute("type", "password")
    if (document.getElementById("password2") != null) {
      document.getElementById("password2").setAttribute("type", "password")
    }
  }
}

function resetCheck(e) {
  document.getElementById("check").innerHTML = '&#8635;';
  document.getElementById("submit").hidden = false;
  document.getElementById("wait").hidden = true;
}

function resetWait(e) {
  document.getElementById("submit").hidden = false;
  document.getElementById("wait").hidden = true;
}

function waitResponse(e) {
  document.getElementById("submit").hidden = true;
  document.getElementById("wait").hidden = false;
}

function check(e){
  em = document.getElementById("error_messages")
  // em2 = document.getElementById("error_messages2")
  // console.log("text content:", em.getAttribute("name"), em.innerHTML, em2.innerHTML)
  // console.log("text content:", em.innerText, em2.innerText)
  if (em.innerText != "") {
    f = document.getElementById("sign_up_form")
    f.removeChild(e.target)
    // s = document.createElement("span")
    // s.innerHTML = `<span 
    //     name="check" id="check" 
    //     class="relative right-[25px] top-[30px] text-sm"
    //   > &#128681;
    //   </span>`
    // f.appendChild(s)
    // e.target.innerHTML = `<span class='text-sm'>&#128681;</span>`
  }
}

function checkCheckError(e){
  setTimeout(function(){
    check(e)
  }, 200);
  // https://stackoverflow.com/qu estions/18749591/encode-html-entities-in-javascript
  // &#8635;
  // const encodedStr = rawStr.replace(/[\u00A0-\u9999<>\&]/g, i => '&#'+i.charCodeAt(0)+';')
  // var encodedStr = e.target.innerText.replace(/[\u00A0-\u9999<>\&\#]/g, function(i) {
    //   return '&#'+i.charCodeAt(0)+';';
  // });
  // console.log(e.target)
  // console.log("check greska", encodedStr )
  // console.log("check greska", e.target.innerText == `&#128681;`)
  // console.log("check greska", e.target.innerHTML == `&#128681;`)
  // console.log("check greska", e.target.innerText == `🚩`)
  // console.log("check greska", e.target.innerHTML == `🚩`)
  // console.log("check greska", encodedStr[0], encodedStr[1], encodedStr[2])
  // if (e.target.innerHTML == `🚩`) {
  //   document.getElementById("error_messages").innerText = "Internal server error, refresh page and try later to check available names."
  // } 
}