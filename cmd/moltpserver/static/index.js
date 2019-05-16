var nf = 0;
function add(){
  nf++;
  let a = document.querySelector("#formulas")
  let e = document.createElement("li")
  e.id = String("id"+nf)

  input = document.createElement("input")
  input.setAttribute("type", "text")
  input.setAttribute("onblur", String(`render('${e.id}')`))

  button = document.createElement("button")
  button.setAttribute("type", "button")
  button.setAttribute("onclick", "rem(this.parentNode)")
  button.innerText = "-"

  div = document.createElement("div")
  div.classList.add("latex")

  e.appendChild(input)
  e.appendChild(button)
  e.appendChild(div)
  a.append(e)
}

function rem(elem){
  elem.remove()
}

function render(eid){
  let i1 = document.querySelector(String(`#${eid} > input:nth-child(1)`))
  // let i2 = document.querySelector('#'+eid+' :nth-child(2)')
  let t  = document.querySelector(String(`#${eid} .latex`))
  // katex.render(String.raw`${i1.value} \gets ${i2.value}`, t);
  katex.render(String.raw`${i1.value}`, t);
}

function solve(){
  var data = []
  document.querySelectorAll('#formulas li').forEach((elem,i)=>{
    data.push({'oid':i, 'left':elem.children[0].value, 'right':elem.children[1].value})
  })

  return fetch("/solve", {
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
        for(var k in data){
          e = document.createElement('div')
          a.appendChild(e)
          katex.render(String.raw(`{\bf${k}:} ${data[k]}`), e);
        }
      })
    }
  })
  .catch(error => console.error(`Fetch Error =\n`, error));
}
