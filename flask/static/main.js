"use strict";

const add_btn = document.getElementById("addButton");
let material_list = [];
if (localStorage.getItem("materials") != null) {
    material_list = (JSON.parse(localStorage.getItem("materials")));
}
console.log(material_list);

add_btn.addEventListener("click", () => {
    let m_name = document.getElementById("material_input").value;
    if (material_list != null) {
        material_list.push(m_name);
        localStorage.setItem('materials', JSON.stringify(material_list));
    } else {
        localStorage.setItem('materials', [m_name]);
    }
})

for (var val in localStorage.getItem("materials")) {
    const table = document.getElementById("table_body");
    const tr = document.createElement("tr");

    const td1 = document.createElement("td");
    const td_input1 = document.createElement("input");
    td_input1.setAttribute("type", "checkbox");
    td_input1.setAttribute("name", "chk1");
    td_input1.setAttribute("value", "1");

    const td2 = document.createElement("td");
    const td_input2 = document.createElement("input");
    td_input2.setAttribute("type", "checkbox");
    td_input2.setAttribute("name", "chk1");
    td_input2.setAttribute("value", "2");

    const td_name = document.createElement("td");
    td_name.innerHTML = val;

    const td_btn = document.createElement("td");

    const btn = document.createElement("button");
    btn.setAttribute("type", "button");
    btn.setAttribute("class", "btn btn-danger");
    btn.innerHTML = "削除";

    td_btn.appendChild(btn);
    td2.appendChild(td_input2);
    td1.appendChild(td_input1);

    tr.appendChild(td1);
    tr.appendChild(td2);
    tr.appendChild(td_name);
    tr.appendChild(td_btn);

    table.appendChild(tr);

}

