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
          for(var k in data){
          e = document.createElement('div')
          a.appendChild(e)
          katex.render(String.raw(`{\bf${k}:} ${data[k]}`), e);
        }
        }
      })
    }
  })
  .catch(error => alert(`Fetch Error =${error}\n`));
}
