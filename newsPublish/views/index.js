
window.onload=function (ev) {

        $(".dels").click(function () {
//  confirm ("是否确认删除?")  不管是否确认删除 它都会删除 触发false
            if(!confirm ("是否确认删除?")){
                //只要返回值为false 就触发删除
                return false
            }
        })
        $("#select").change(function () {
            $("#form").submit()
        })
    }

