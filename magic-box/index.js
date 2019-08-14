//魔方单元对应标识
var rubikSide = {
    up : [1, 2, 3, 4, 5, 6, 7, 8, 9],
    bottom : [18, 19, 20, 21, 22, 23, 24, 25, 26],
    left : [7, 4, 1, 15, 13, 10, 24, 21, 18],
    right : [9, 6, 3, 17, 14, 12, 26, 23, 20],
    front : [7, 8, 9, 15, 16, 17, 24, 25, 26],
    back : [1, 2, 3, 10, 11, 12, 18, 19, 20]
}

var interface = [-30,-30,0]

var lock = {
    isLock: false
}

function creatColor(){
    var obj = {};
    for(let i = 1; i < 27; i++){
        obj[`${i}`] = {
            up: 'size-color-red',
            bottom: 'size-color-orange',
            left: 'size-color-blue',
            right: 'size-color-green',
            front: 'size-color-wihte',
            back: 'size-color-yellow'
        }
    }
    return obj;
}

var modulesColor = creatColor();

var keySide = {
    115 : ["left", 1],//1代表正旋-1代表反旋
    100 : ["left", -1],
    102 : ["front", 1],
    103 : ["front", -1],
    104 : ["back", 1],
    106 : ["back", -1],
    116 : ["up", 1],
    121 : ["up", -1],
    118 : ["bottom", 1],
    98 : ["bottom", -1],
    107 : ["right", 1],
    108 : ["right", -1]
};

//设置魔方模型
function setPosition(num){
    var tar = $(`.cube${num}>.sigle-module`)[0];
    var matrixTrans = [100,100,0];
    if(arrTest(num, rubikSide.up)){
        matrixTrans[1] -= 100;
    }
    if(arrTest(num, rubikSide.bottom)){
        matrixTrans[1] += 100;
    }
    if(arrTest(num, rubikSide.left)){
        matrixTrans[0] -= 100;
    }
    if(arrTest(num, rubikSide.right)){
        matrixTrans[0] += 100;
    }
    if(arrTest(num, rubikSide.front)){
        matrixTrans[2] += 100;
    }
    if(arrTest(num, rubikSide.back)){
        matrixTrans[2] -= 100;
    }
    tar.style.transform = `matrix3d(1,0,0,0,0,1,0,0,0,0,1,0,${matrixTrans[0]},${matrixTrans[1]},${matrixTrans[2]},1)`;
}

//获取matrix数组
function getMatrix(str){
    var reg = /(?=\()([\d\D])+(?=\))/ig;
    if(reg.test(str)){
        return str.match(reg)[0].split(',');
    } 
}

//检测数组中是否含有某个数值
function arrTest(num, arr){
    for(let i in arr){
        if(arr[i] == num){
            return true;
        }
    }
    return false;
}


function setCubePos(){
    //获取所有旋转模块
    var cubes = $('.rotate-module');
    var cubesL = cubes.length;
    //循环模块
    for(let i = 1; i <= cubesL; i++){
        setPosition(i);
    }
    
}


setCubePos();

//设置旋转模块
function rotateModule(){
    document.onkeypress = function(e){
        if(keySide[`${e.keyCode}`]){
            if(!lock.isLock){
                lock.isLock = true;
                var rotateType = keySide[`${e.keyCode}`];
                rotate(rotateType[0], rotateType[1]);
            }
        }
    }
}



rotateModule()

//封装魔方旋转函数
function rotate(type, dir){ 
    var size = $('.rotate-size')
    var turn = 90*dir;
    var axile;
    if(type == 'up' || type == 'bottom'){
        axile = [0, turn, 0];
    }else if(type == 'left' || type == 'right'){
        axile = [turn, 0, 0];
    }else{
        axile = [0, 0, turn];
    }
    var rotateArr = rubikSide[type];
    var len = rotateArr.length;
    for(let i = 0; i < len; i++){
        var target = $(`.cube${rotateArr[i]}`)[0];
        size.append(target);      
    }
    size[0].style.transition = 'all 0.4s';
    size[0].style.transform = `rotateX(${axile[0]}deg) rotateY(${axile[1]}deg) rotateZ(${axile[2]}deg)`;
    size[0].addEventListener("transitionend", function transEnd(){
        if(count.steps == 1){
            stepCount();
  
        }
        lock.isLock = false;
        this.style.transition = '';
        for(let i in rotateArr){
            colorChange(rotateArr[i], type, dir);
            if(i == len-1){
                if(type == 'up' || type == 'bottom'){
                    changeBox(type, -1*dir)
                }else{
                    changeBox(type, dir);
                }
            }
        }
        
        this.style.transform = `rotateX(0deg) rotateY(0deg) rotateZ(0deg)`;
        size.after($('.rotate-size .rotate-module'));
        size.empty();
        this.removeEventListener("transitionend", transEnd);
    });

}




function colorChange(id, type, dir){
    var target = modulesColor[id];
    var orign = {};
    for(let key in target){
        orign[key] = target[key];
    }
    var changeSize, temp;
    if(type == 'up' || type == 'bottom'){
        changeSize = ['front','right','back','left'];
    }else if(type == 'left' || type == 'right'){
        changeSize = ['front','up','back','bottom'];
    }else{
        changeSize = ['left','up','right','bottom'];
    }
    if(dir === 1){
        var temp = target[changeSize[3]];
        target[changeSize[3]] = target[changeSize[2]];
        target[changeSize[2]] = target[changeSize[1]];
        target[changeSize[1]] = target[changeSize[0]];
        target[changeSize[0]] = temp;
    }else{
        var temp = target[changeSize[0]];
        target[changeSize[0]] = target[changeSize[1]];
        target[changeSize[1]] = target[changeSize[2]];
        target[changeSize[2]] = target[changeSize[3]];
        target[changeSize[3]] = temp;
    }
    for(let i = 0; i < changeSize.length; i++){
        $(`.cube${id} .sigle-${changeSize[i]}`).removeClass(orign[changeSize[i]]);
        $(`.cube${id} .sigle-${changeSize[i]}`).addClass(modulesColor[id][changeSize[i]]);
    }
    
}

function changeBox(type ,dir){
    var targets = rubikSide[type],
        len = targets.length;
    var allColor = 'size-color-red size-color-orange size-color-blue size-color-green size-color-wihte size-color-yellow'
    if(dir === 1){
        exchange(modulesColor, targets[0], targets[6]);
        exchange(modulesColor, targets[6], targets[8]);
        exchange(modulesColor, targets[8], targets[2]);
        exchange(modulesColor, targets[3], targets[7]);
        exchange(modulesColor, targets[7], targets[5]);
        exchange(modulesColor, targets[5], targets[1]); 
    }else{
        exchange(modulesColor, targets[8], targets[2]);
        exchange(modulesColor, targets[6], targets[8]);
        exchange(modulesColor, targets[0], targets[6]);
        exchange(modulesColor, targets[5], targets[1]);
        exchange(modulesColor, targets[7], targets[5]);
        exchange(modulesColor, targets[3], targets[7]);
    }
    for(let i = 0; i < len; i++){
        for(let key in modulesColor[targets[i]]){
            $(`.cube${targets[i]} .sigle-${key}`).removeClass(allColor);
            $(`.cube${targets[i]} .sigle-${key}`).addClass(modulesColor[targets[i]][key])
        }
    }
}

function exchange(tar, num1, num2){
    var temp = {};
    for(let key in tar[num1]){
        temp[key] = tar[num1][key];
    }
    for(let key in tar[num2]){
        tar[num1][key] = tar[num2][key];
    }
    for(let key in temp){
        tar[num2][key] = temp[key];
    }
}



function interfaceChange(){
    var doc = document;
    doc.addEventListener('keypress', function(e){
        var key = e.keyCode;
        if(key == '120' || key == '99'){
            key == 120 ? interface[1] -= 2 : interface[1] += 2;
        }else if(key == '110' || key == '109'){
            key == 110 ? interface[0] += 2 : interface[0] -= 2;
        }
        $('.box-content')[0].style.transform = `rotateX(${interface[0]}deg) rotateY(${interface[1]}deg) rotateZ(0deg)`;
    })




}
interfaceChange()

function startCount(){

}

function stopCount(){

}

var count = {
    lock: 1,
    lockStep: 1,
    time: 0,
    step : 0,
    steps : 0
}

function stepCount(){
    count.steps ++;
    $('.step')[0].innerHTML = count.steps;
}

function upsetLayou(){
    if(count.lockStep == 1){
        count.lockStep = 0;
        autoMove()
        var timer = setInterval(() => {
            if(count.step < 10){
                autoMove()
            }else{
                count.lockStep = 1;
                count.step = 0;
                clearInterval(timer);
                return false
            }
        }, 480);
    }
    
    
}


function autoMove(){
    var options = [115,100,102,103,104,106,116,121,118,98,107,108];
    count.step++;
    var random = parseInt(Math.random()*3) + parseInt(Math.random()*9);
    var rotateType = keySide[options[random]];
    rotate(rotateType[0], rotateType[1]);
}


function eventClick(){
    var timer1;
    
    $('.start-count')[0].onclick = function (e){
        $('.time')[0].innerHTML = count.time;
        if(count.lock == 1){
            count.lock =0;
            timer1 = setInterval(() => {
                count.time ++;
                $('.time')[0].innerHTML = count.time;
                }, 1000);
        }
        
    }
    $('.stop-count')[0].onclick = function (e){
        clearInterval(timer1)
        count.time = 0;
        count.lock = 1;
    }
    
}

eventClick()