function showClasses() {
    document.getElementById("searchDropdown").classList.toggle("show")
}

function filterFunction() {
    let input, filter, ul, li, a, i
    input = document.getElementById("searchInput")
    filter = input.value.toUpperCase()
    div = document.getElementById("searchDropdown")
    a = div.getElementsByTagName("a")
    for (i = 0; i < a.length; i++) {
        txtValue = a[i].textContent || a[i].innerText
        if (txtValue.toUpperCase().indexOf(filter) > -1) {
            a[i].style.display = ""
        } else {
            a[i].style.display = "none"
        }
    }
}