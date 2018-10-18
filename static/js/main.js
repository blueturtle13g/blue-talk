
    jQuery(document).ready(function($){

        $("#upPro").click(function (e) {
            // var selectedCountry = $("#private option:selected").val();
            // alert($("#private option:selected").val());
            e.preventDefault();
            $.ajax({
                url: "/profile/{{.User.Id}}",
                type: 'post',
                dataType: 'html',
                data : { username : $("#username").val(),
                            firstName : $("#firstName").val(),
                            lastName : $("#lastName").val(),
                            phone : $("#phone").val(),
                            email : $("#email").val(),
                            quote : $("#quote").val(),
                            pri : $("#private option:selected").val()
                },
                success : function(data) {
                    if (data !== "done"){
                        errAlert(data)
                    }else{
                        window.location.reload();
                    }
                }
            });
        });


        $("#delPP").click(function (e) {
            e.preventDefault();
            console.log("got to delpp")
            let confirmBox = $("#confirmBox");
            confirmBox.find(".message").html("<h5> Are you sure?</h5>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/profile/{{.User.Id}}",
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "DelPP",
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                        } else {
                            errAlert(data)
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });

        $("#upPass").click(function (e) {
            e.preventDefault();
            $.ajax({
                url: "/profile/{{.User.Id}}",
                type: 'post',
                dataType: 'html',
                data : { submit : "UpPass",
                        cPass : $("#cPass").val(),
                        newPass : $("#newPass").val(),
                        confirmPass : $("#confirmPass").val(),
                },
                success : function(data) {
                    if (data !== "done"){
                        errAlert(data)
                    }else{
                        congAlert("Your Password has been updated successfully.")
                    }
                }
            });
        });

        uploadImage();
        function uploadImage() {
            var button = $('.images .pic');
            var uploader = $("#picsInput");
            var images = $('.images');

            button.on('click', function () {
                if ($(".postImgHelper").length >= 6){
                    errAlert("You can't uer more than 6 pictures in a post.")
                }else{
                    uploader.click()
                }
            });

            uploader.on('change', function () {
                var reader = new FileReader();
                reader.onload = function(event) {
                    images.prepend('<div class="postImgHelper" style="background-image: url(\'' + event.target.result + '\');" rel="'+ event.target.result  +'"><span>remove</span></div>')
                };
                reader.readAsDataURL(uploader[0].files[0])

            });

            images.on('click', '.postImgHelper', function () {
                $(this).remove()
            })

        }

        $("#shP").click(function() {
            let passes = $(".passes");
            if (passes.attr('type') === "password") {
                passes.attr('type', 'text');
            } else {
                passes.attr('type', 'password');
            }
        });

        // proFields
        $(".proFields").click(function () {
            $(".proFields").removeClass("proFieldsActive");
            $(this).addClass("proFieldsActive");
        });

        // for file uploading
        var readURL = function(input) {
            console.log("readUrl");
            if (input.files && input.files[0]) {
                var reader = new FileReader();

                reader.onload = function (e) {
                    $('.avatar').attr('src', e.target.result);
                };

                reader.readAsDataURL(input.files[0]);
            }
        };

        $("#upload").on('change', function(){
            readURL(this);
        });

        // for alerts

        $('body').click(function(e) {
            if (!$(e.target).closest('.alertBox').length){
                $(".alertBox").hide();
            }
        });

        function errAlert(msg) {
            let errBox = $("#errBox");
            errBox.find(".message").html(msg);
            errBox.find(".no").unbind().click(function () {
                $("#confirmBox").hide();
                errBox.hide();
            });
            errBox.show("slide", {direction: "up"}, 500);
        }

        function congAlert(msg) {
            let congBox = $("#congBox");
            congBox.find("div").html(msg);
            congBox.find("button").unbind().click(function () {
                window.location.reload(true);
            });
            congBox.show("slide", {direction: "up"}, 500);
        }

        // for commenting ajaxes
        $("#sendCom").click(function (e) {
            e.preventDefault();
            let storyId = $("#sendCom").data("storyid");
            $.ajax({
                url: "/story/"+storyId,
                type: 'post',
                dataType: 'html',
                data: {
                    submit: "com",
                    text: $("#comText").val()
                },
                success: function (data) {
                    if (data === "done") {
                        window.location.reload();
                    } else {
                        errAlert(data);
                    }
                }
            });
        });
        function indexInClass(node) {
            var collection = document.getElementsByClassName(node.className);
            for (var i = 0; i < collection.length; i++) {
                if (collection[i] === node)
                    return i;
            }
            return -1;
        }
        OpenIndex();
        function OpenIndex(){
            var index = window.sessionStorage.getItem("openRes");
            $(".ResField").eq(index).show();
        }

        $(".exCom").click(function() {
            window.sessionStorage.setItem("openRes", indexInClass(this));
            $(this).parent().next().toggle("slide", { direction: "up" }, 500);
        });

        $(".exEdit").click(function(e) {
            e.preventDefault();
            var textBox = $("#textBox");
            var comId = $(this).data("comid");
            textBox.find(".message").html("<h5>Edit Your Comment</h5>");
            textBox.find("textarea").val($(this).data("comtext"));
            textBox.find(".send").unbind().click(function () {
                $.ajax({
                    url: "/story/"+$(this).data("storyid"),
                    type: 'post',
                    dataType: 'html',
                    data: {
                        submit: "EditCom",
                        text: textBox.find("textarea").val(),
                        comId: comId,
                    },
                    success: function (data) {
                        if (data === "done") {
                            window.location.reload();
                        } else {
                            errAlert(data);
                        }
                    }
                });
            });
            textBox.find(".cancel").unbind().click(function () {
                textBox.hide();
                $("#errBox").hide();
            });
            textBox.show("slide", {direction: "up"}, 500);
            textBox.find("textarea").focus();
        });

        $(".sendResCom").click(function (e) {
            e.preventDefault();
            let comId = $(this).data("comid");
            let storyId = $(this).data("storyid");
            let text = $(this).parent().parent().children("textarea").val();
            $.ajax({
                url: "/story/"+storyId,
                type: 'post',
                dataType: 'html',
                data: {
                    submit: "resCom",
                    text: text,
                    comId: comId,
                },
                success: function (data) {
                    if (data === "done") {
                        window.location.reload();
                    } else {
                        errAlert(data);
                    }
                }
            });
        });
        // end of commenting ajaxes
        $("#outSure").click(function (e) {
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            confirmBox.find(".message").html("<h5>Are you sure?</h5>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: '/out',
                    type: 'get',
                    dataType: 'html',
                    success: function (data) {
                        if (data === "done") {
                            window.location.reload()
                        } else {
                            errAlert("<h5>Your operation doesn't get completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });

        $("#dImgPSure").click(function (e) {
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            confirmBox.find(".message").html("<h5>Are you sure?</h5>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/profile/{{.User.Id}}/edit",
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "DeleteImg",
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                        } else {
                            errAlert("<h5>Your operation doesn't get completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });

        $("#dImgPSure").click(function (e) {
            var storyId = $(this).data("storyid");
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            confirmBox.find(".message").html("<h5>Are you sure?</h5>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/story/"+storyId+"/edit",
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "DeleteImg",
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                        } else {
                            errAlert("<h5>Your operation doesn't get completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });

        $("#dpSure").click(function (e) {
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            confirmBox.find(".message").html("<h5> Are you sure you want to delete your account?</h5><p>It can't be retrieved later</p>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/profile/{{.User.Id}}/edit",
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "Delete",
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                            window.location = "/";
                        } else {
                            errAlert("<h5>Your operation doesn't be completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });
        $("#dpSure").click(function (e) {
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            let storyid = $("#dsSure").data("storyid");
            console.log(storyid);
            confirmBox.find(".message").html("<h5> Are you sure you want to delete this story?</h5><p>It can't be retrieved later</p>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/story/"+storyid+"/edit",
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "Delete",
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                            window.location = "/";
                        } else {
                            errAlert("</h5>Your operation doesn't be completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });
        $(".drSure").click(function (e) {
            e.preventDefault();
            let confirmBox = $("#confirmBox");
            let comId = $(this).data("comid");
            let storyId = $(this).data("storyid");
            confirmBox.find(".message").html("<h5>Are you sure?</h5>");
            confirmBox.find(".yes").unbind().click(function () {
                $.ajax({
                    url: "/story/"+storyId,
                    type: 'post',
                    dataType: 'html',
                    data : { submit : "DeleteCom",
                        comId : comId
                    },
                    success : function(data) {
                        if (data === "done"){
                            window.location.reload();
                        } else {
                            errAlert("</h5>Your operation doesn't be completed</h5>")
                        }
                    }
                });
            });
            confirmBox.find(".no").unbind().click(function () {
                confirmBox.hide();
            });
            confirmBox.show("slide", {direction: "up"}, 500);
        });
        // end of alerts

        //validation form
        $('input:submit').on('click', function(e) {
            if ($(".inputErr").is(":visible")){
                e.preventDefault();
                $('#submitWarn').show();
            }
        });
        $('#username').on('blur', function(){
            let username = $('#username').val();
            $('#submitWarn').hide();
            if (username === "") {
                $('#usernameWarn').hide();
                return;
            }else if (username.length < 3){
                $('#usernameWarn').show();
                $("#usernameWarn").addClass("inputErr").removeClass("inputOk");
                $('#usernameWarn').html("<strong>At Least 3 Characters</strong>");
                return;
            }
            $.ajax({
                url: '/register',
                type: 'post',
                dataType: 'html',
                data : { submit: "CheckName",
                    writerId: $("#username").data("userrid"),
                    username: $("#username").val(),
                },
                success : function(data) {
                    if (data === "Invalid"){
                        $("#usernameWarn").addClass("inputErr").removeClass('inputOk');
                        $("#usernameWarn").html("<strong title='It cannot have space between and start with numbers, it also should be in english.'>Invalid Name</strong>");
                        $("#usernameWarn").show();
                    }else if (data === "Taken") {
                        $("#usernameWarn").addClass("inputErr").removeClass('inputOk');
                        $("#usernameWarn").html("<strong title='some user has already chosen this name, so it is not available to you.'>Already Taken</strong>");
                        $("#usernameWarn").show();
                    } else if (data === "Available") {
                        $("#usernameWarn").addClass("inputOk").removeClass("inputErr");
                        $("#usernameWarn").html("<strong title='No user has chosen this name, so it is available to you.'>Available</strong>");
                        $("#usernameWarn").show();
                    } else {
                        $("#usernameWarn").hide();
                    }
                }
            });
        });

        $('#email').on('blur', function(){
            $('#submitWarn').hide();
            var email = $("#email").val();
            if (email === ""){
                $('#emailWarn').hide();
            }else if (validateEmail(email)) {
                $('#emailWarn').hide();
            }else{
                $('#emailWarn').show();
            }
        });
        $('#password').on('blur', function(){
            $('#submitWarn').hide();
            let password = $('#password').val();
            if (password === ""){
                $('#passWarn').hide();
            }else if (password.length < 8) {
                $('#passWarn').show();
            }else{
                $('#passWarn').hide();
            }
        });
        function validateEmail(email) {
            var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
        $('#confirm').on('blur', function(){
            $('#submitWarn').hide();
            let password = $('#password').val();
            let confirm = $('#confirm').val();
            if (confirm === "") {
                $('#confWarn').hide();
            }else if (confirm !== password){
                $('#confWarn').show();
            }else{
                $('#confWarn').hide();
            }
        });
        // end of validation form
    });