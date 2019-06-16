solution = document.querySelector('#solution')

function render(formula, renderer){
  let input = document.querySelector(String(`#${formula}`))
  let where = document.querySelector(String(`#${renderer}`))
  katex.render(String(`${input.value}`), where)
  if (where.childElementCount < 1) {
    alert("Error rendering your formula input. It will not be submitted. Check your input.")
  }
}

function fillSolution(data) {
  var k = 0
  var s = data[k]
  while( s != undefined ){
    li = document.createElement('li')

    solution.appendChild(li)
    d1 = document.createElement('div')
    d1.classList.add("sequntdivider")
    li.appendChild(d1)
    d2 = document.createElement('div')
    d2.classList.add("sequntsegment")
    li.appendChild(d2)
    d3 = document.createElement('div')
    d3.classList.add("sequntdivider")
    li.appendChild(d3)
    d4 = document.createElement('div')
    d4.classList.add("sequntsegment")
    li.appendChild(d4)
    d5 = document.createElement('div')
    d5.classList.add("sequntdivider")
    li.appendChild(d5)

    d1.innerText = String(`${s["name"]}`)
    katex.render(String(`${s["left"]}`), d2);
    katex.render("\\leftarrow", d3);
    katex.render(String(`${s["right"]}`), d4);
    d5.innerText = String(`${s["just"]}`)
    k++
    s = data[k]
  }
}

function prove(){
  var data = {'oid':0, 'formula':document.querySelector("#f1").value}
  solution.innerHTML = ''
  document.querySelector('#soltitle').innerText = "Solution"

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
        alert(String(`${data["info"])}`)
        if (response.status == 500) {
          if (data != null && data != "null")  {
            document.querySelector('#soltitle').innerText = "Partial result"
            fillSolution(data["result"])
          }
        }
      })
    } else {
      response.json().then(function(data){
        if (data == null || data == "null")  {
          alert("Empty reponse!")
        } else {
          fillSolution(data)
        }
      })
    }
  })
  .catch(error => alert(`Fetch Error =${error}\n`));
}
