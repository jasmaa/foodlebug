
// convert rating to stars
function numToStars(n, target){
    n = parseInt(n);
    res = "";
    for(let i = 0; i < 5; i++){
        if(i < n){
            res += `<span class="fas fa-star"></span>\n`;
        }
        else{
            res += `<span class="far fa-star"></span>\n`;
        }
    }
    document.getElementById(target).innerHTML = res;
}