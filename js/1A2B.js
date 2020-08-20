$(document).ready(function () 
{
    $('#submit').click(function (){
        if(check()){
            $.ajax({
                url: 'http://127.0.0.1:8085/game',
                cache: false,
                dataType: 'html',
                type:'GET',
                xhrFields: {
                    withCredentials: true
                },
                crossDomain:true,       //允許跨域
                data: { answer: $('#answer').val()},
                error: function(xhr) {
                alert('Ajax request 發生錯誤');
            },
            success: function(response) {
                $('#msg').html(response);
                $('#msg').fadeIn();
            }
            });
        } 
    });
})
function check() 
{
    var answer =  $("#answer").val();
    if(answer.length !== 4){
        alert("未輸滿4個數字，請填滿，謝謝");
        eval("document.form1['answer'].focus()");
        return false;
    }else{
        return true;
    }
}
