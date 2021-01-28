"use strict";

const refreshTable = () => {
    const table = document.getElementById("table_body");
    table.innerHTML = "";

    material_list.forEach(elem => {
        const tr = document.createElement("tr");

        const td1 = document.createElement("td");
        td1.setAttribute("id", "td_checkbox");
        const td_input1 = document.createElement("input");
        td_input1.setAttribute("type", "checkbox");
        td_input1.setAttribute("name", "chk1");
        // td_input1.setAttribute("value", "1");
        td_input1.setAttribute("value", elem);
        td_input1.setAttribute("id", "check_box");
        td_input1.setAttribute("class", "align-middle");

        const td2 = document.createElement("td");
        td2.setAttribute("id", "td_checkbox");
        const td_input2 = document.createElement("input");
        td_input2.setAttribute("type", "checkbox");
        td_input2.setAttribute("name", "chk2");
        // td_input2.setAttribute("value", "2");
        td_input2.setAttribute("value", elem);
        td_input2.setAttribute("id", "check_box");
        td_input2.setAttribute("class", "align-middle");


        const td_name = document.createElement("td");
        td_name.setAttribute("id", "material_name");
        td_name.innerHTML = elem;

        const td_btn = document.createElement("td");

        const btn = document.createElement("button");
        btn.setAttribute("type", "button");
        btn.setAttribute("class", "btn btn-danger");
        btn.innerHTML = "削除";
        btn.addEventListener('click', function () {
            const target_tr = this.parentNode.parentNode;
            const target_name = target_tr.children[2].innerHTML;
            material_list.splice(material_list.indexOf(target_name), 1); // 配列から削除
            localStorage.setItem('materials', JSON.stringify(material_list)); // localStorageを更新
            refreshButton(); // ボタンの更新
            target_tr.parentNode.deleteRow(tr.sectionRowIndex); // 行の削除
        });
        td_btn.appendChild(btn);
        td2.appendChild(td_input2);
        td1.appendChild(td_input1);

        tr.appendChild(td1);
        tr.appendChild(td2);
        tr.appendChild(td_name);
        tr.appendChild(td_btn);

        table.appendChild(tr);
    });
};

const refreshButton = () => {
    const empty_message = document.getElementById("empty_message");
    const search_btn = document.getElementById("search_btn");
    if (material_list.length == 0) {
        search_btn.style.visibility = "hidden";
        empty_message.style.visibility = "visible";
    } else {
        search_btn.style.visibility = "visible";
        empty_message.style.visibility = "hidden";
    }
};

const add_btn = document.getElementById("addButton");
let material_list = [];
if (localStorage.getItem("materials") != null) {
    material_list = (JSON.parse(localStorage.getItem("materials")));
}
console.log(material_list);
refreshButton();
refreshTable();

add_btn.addEventListener("click", () => {
    let m_name = document.getElementById("material_input").value;
    const warning_msg = document.getElementById("warning_msg");

    if (m_name != "") {
        if (material_list != null) {
            material_list.push(m_name);
            localStorage.setItem('materials', JSON.stringify(material_list));
        } else {
            localStorage.setItem('materials', [m_name]);
        }
        warning_msg.style.visibility = "hidden"; // 警告メッセージの非表示
        refreshTable();
        refreshButton();
    } else {
        warning_msg.style.visibility = "visible"; // 警告メッセージの表示
    }
});

// 検索処理
const search_btn = document.getElementById("search_btn");
search_btn.addEventListener("click", () => {
    const ignore_list = [];
    const target_list = [];
    const chk1 = document.form1.chk1;
    const chk2 = document.form1.chk2;
    console.log(chk1);

    if (chk1.length) { // データが2件以上のときの処理
        for (let i = 0; i < chk1.length; i++) {
            if (chk1[i].checked) {
                ignore_list.push(chk1[i].value);
            }
            if (chk2[i].checked) {
                target_list.push(chk1[i].value);
            }
        }
    } else { // データが1件のときの処理
        if (chk1.checked) {
            ignore_list.push(chk1.value);
        }
        if (chk2.checked) {
            target_list.push(chk1.value);
        }
    }

    console.log(ignore_list);
    console.log(target_list);
    if(ignore_list.length > 0 && target_list.length > 0){
        window.location.href = `/result?name=${target_list}&ignore=${ignore_list}`;
    } else if (ignore_list.length == 0 && target_list.length > 0){
        window.location.href = `/result?name=${target_list}`;
    } else {
        // 検索対象未選択
    }
})

jQuery(function ($) {
    $('input:checkbox').click(function () {
        $(this).closest('.table tr').find('input:checkbox').not(this).prop('checked', false);
    });
});