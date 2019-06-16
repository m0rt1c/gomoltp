function render(formula, renderer){
  let input = document.querySelector(String(`#${formula}`))
  let where = document.querySelector(String(`#${renderer}`))
  katex.render(String(`${input.value}`), where)
  if (where.childElementCount < 1) {
    alert("Error rendering your formula input. It will not be submitted. Check your input.")
  }
}

function prove(){
  var data = {'oid':0, 'formula':document.querySelector("#f1").value}

  return fetch("/prover", {
    method: "POST",
    headers: {
      "Content-Type": "application/json; charset=utf-8",
    },
    body: JSON.stringify(data),
  })
  .then(function(response) {
    if (response.status != 200) {
      response.json().then(function(data){
        alert("Info: "+String(data["info"]))
      })
    } else {
      response.json().then(function(data){
        a = document.querySelector('#solution')
        a.innerHTML = ''
        if (data == null || data == "null")  {
          alert("Empty reponse!")
        } else {
          var k = 0
          var s = data[k]
          while( s != undefined ){
            li = document.createElement('li')
            a.appendChild(li)
            d1 = document.createElement('div')
            d1.classList.add("sequntsegment")
            d2 = document.createElement('div')
            d2.classList.add("sequntdivider")
            d3 = document.createElement('div')
            d3.classList.add("sequntsegment")
            d4 = document.createElement('div')
            d4.classList.add("sequntsegment")
            d5 = document.createElement('div')
            d5.classList.add("sequntsegment")
            li.appendChild(d5)
            li.appendChild(d1)
            li.appendChild(d2)
            li.appendChild(d3)
            li.appendChild(d4)
            d5.innerText = String(`${s["name"]}`)
            katex.render(String(`${s["left"]}`), d1);
            katex.render("\\leftarrow", d2);
            katex.render(String(`${s["right"]}`), d3);
            d4.innerText = String(`${s["just"]}`)
            k++
            s = data[k]
          }
        }
      })
    }
  })
  .catch(error => alert(`Fetch Error =${error}\n`));
}
