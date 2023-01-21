$(document).ready(function () {
    $('#logout').submit(function (e) {
        Cookies.remove('auth-session');
    });
});