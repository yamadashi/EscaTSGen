var isBarTime = false;

//Webフォントのプリロード
WebFont.load({
    custom: {
        families: ['BlackChancery'],
        urls: ['../css/style.css']
    },
    loading: function () {
        console.log('loading');
    },
    active: function () {
        console.log('active');
    },
    inactive: function () {
        console.log('inactive');
    }
});

// Canvasの準備
var canvas = document.getElementById('canvas');
var ctx = canvas.getContext('2d');
canvas.style.display = "none"; //canvasから画像を生成するがcanvas自体は表示しない


//フォーム追加
var timeline = document.getElementById("timeline");
const input_form = document.getElementById("input_form");
function clone_form() {
    var new_form = input_form.cloneNode(true);
    new_form.style.display = 'inline';
    //フォーム消去イベント
    var delete_button = getElementsWithAttribute(new_form, "name", "delete")[0];
    delete_button.addEventListener('click', function () {
        timeline.removeChild(new_form);
    });

    timeline.appendChild(new_form);
}
var add_button = document.getElementById("add_button");
add_button.addEventListener('click', clone_form);
//一つ目のフォーム
clone_form();


//画像情報
var baseImageInfo = null;
function initBaseImageInfo() {
    if (isBarTime) {
        baseImageInfo = {
            size: {
                width: 480,
                height: 640
            },
            timelineArea: {
                top: 212,
                bottom: 437,
                leftOffset: 70
            },
            //幅が均一でないためテーブルで位置を保持
            timePosTable: new Map(
                [
                    [22, 212],
                    [23, 240],
                    [0, 270],
                    [1, 297],
                    [2, 326],
                    [3, 355],
                    [4, 383],
                    [5, 410],
                    [6, 437]
                ]
            )
        }
    }
    else {
        baseImageInfo = {
            size: {
                width: 800,
                height: 960
            },
            timelineArea: {
                top: 195,
                bottom: 677,
                leftOffset: 100
            },
            //幅が均一でないためテーブルで位置を保持
            timePosTable: new Map(
                [
                    [12, 195],
                    [13, 237],
                    [14, 280],
                    [15, 318],
                    [16, 360],
                    [17, 400],
                    [18, 440],
                    [19, 480],
                    [20, 519],
                    [21, 559],
                    [22, 599],
                    [23, 639],
                    [24, 677]
                ]
            )
        }
    }
}


//画像生成
var gen_button = document.getElementById('gen_button');
function generateImage() {
    // Canvas上に画像を表示
    var img = new Image();
    img.crossOrigin = 'Anonymous';
    img.onload = function () {
        //isBarTimeに応じて背景画像情報の初期化
        initBaseImageInfo();
        //描画準備
        canvas.width = img.width;
        canvas.height = img.height;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        //日付を表示
        putDate();
        //名前とラインを表示
        var shiftInfo = getShiftInfo();
        var eachWidth = (baseImageInfo.size.width - baseImageInfo.timelineArea.leftOffset) / shiftInfo.length;
        for (var i = 0; i < shiftInfo.length; i++) {
            var xpos = baseImageInfo.timelineArea.leftOffset + eachWidth * (i + 1 / 2);
            putNameAndLine(xpos, shiftInfo[i], shiftInfo.length);
        }
        //Canvasを画像として出力
        var data = canvas.toDataURL();
        var outputImg = document.createElement('img');
        outputImg.src = data;
        clearImage();
        result.appendChild(outputImg);
    }
    isBarTime = document.getElementById('barTime').checked;
    img.src = isBarTime ? "base_bar_time.jpg" : "base.jpg";
}
gen_button.addEventListener('click', generateImage, false);


//画像をクリア
var clear_button = document.getElementById("clear_button");
function clearImage() {
    var result = document.getElementById('result');
    if (result.firstChild) result.removeChild(result.firstChild);
}
clear_button.addEventListener('click', clearImage, false);


//日付を表示
function putDate() {
    const date = document.getElementById("date").value.replace(/-/g, '/');
    const subtext = document.getElementById("subtext").value;
    //フォントサイズ
    const dateFontSize = isBarTime ? 25 : (subtext != "" ? 45 : 55);
    const textFontSize = dateFontSize * 0.85;
    const lineHeight = dateFontSize * 1.2;
    //フォントの決定
    //var fontName = [...document.getElementsByName("date-font")].filter(el => el.checked)[0].value; //書き方かっこいいから残しておきたい
    var dateFontName = document.getElementById('date_font').checked ? "BlackChancery" : "ヒラギノ明朝 W6";
    var textFontName = document.getElementById('subtext_font').checked ? "BlackChancery" : "ヒラギノ明朝 W6";
    //contextの指定
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillStyle = '#ffffff';

    if (isBarTime) {
        const textPos = {
            x: baseImageInfo.size.width / 4,
            y: baseImageInfo.size.height / 13
        };
        const datePos = {
            ...textPos,
            x:7*textPos.x/8
        };
        const dateFontSize = 28;
        const lineHeight = dateFontSize * 1.2;
        ctx.font = dateFontSize + "px '" + dateFontName + "'";
        ctx.fillText(date, datePos.x, datePos.y - lineHeight * (subtext != "") / 2);
        const textFontSize = 21;
        ctx.font = textFontSize + "px '" + textFontName + "'";
        ctx.fillText(subtext, textPos.x, textPos.y + lineHeight * (subtext != "") / 2);
    } else {
        var pos = {
            x: baseImageInfo.size.width / 2,
            y: baseImageInfo.size.height / 11
        };
        const dateFontSize = subtext != "" ? 45 : 55;
        const lineHeight = dateFontSize * 1.2;
        ctx.font = dateFontSize + "px '" + dateFontName + "'";
        ctx.fillText(date, pos.x, pos.y - lineHeight * (subtext != "") / 2);
        const textFontSize = dateFontSize * 0.85;
        ctx.font = textFontSize + "px '" + textFontName + "'";
        ctx.fillText(subtext, pos.x, pos.y + lineHeight * (subtext != "") / 2);
    }
}

//入力フォームからシフト情報を受け取る
function getShiftInfo() {
    var shiftInfo = [];
    var elm = timeline.firstElementChild;
    while (elm) {
        var info = {
            name: getElementsWithAttribute(elm, "name", "cast_name")[0].value,
            begin: getElementsWithAttribute(elm, "name", "begin")[0].value,
            end: getElementsWithAttribute(elm, "name", "end")[0].value
        }
        shiftInfo.push(info);
        elm = elm.nextElementSibling;
    }
    return shiftInfo;
}

function getElementsWithAttribute(parent, attributeName, value) {
    var elements = [];
    var elm = parent.firstElementChild;
    while (elm) {
        if (elm.getAttribute(attributeName) === value) {
            elements.push(elm);
        }
        elm = elm.nextElementSibling;
    }
    return elements;
}

//名前とタイムラインを表示
function putNameAndLine(xpos, castInfo, number) {
    //名前
    const fontSize = (isBarTime ? 20 : 27) - (number > 5 ? 2 * (number - 5) : 0);
    ctx.font = fontSize + "px 'ヒラギノ明朝 W6'";
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    ctx.fillStyle = isBarTime ? '#ffffff' : '#375b8b';
    ctx.fillText(castInfo.name, xpos, baseImageInfo.timelineArea.top);
    //タイムライン
    const barWidth = (baseImageInfo.size.width - baseImageInfo.timelineArea.leftOffset) / (isBarTime ? 25 : 30);
    const barTop = calcPos(castInfo.begin);
    const barBottom = calcPos(castInfo.end);
    ctx.fillStyle = '#01aed9';
    ctx.fillRect(xpos - barWidth / 2, barTop, barWidth, barBottom - barTop);
}

function calcPos(timeStr) {
    //数値に変換して範囲外の場合は補正する(例外をどう処理するかは要検討)
    const time = timeStr.split(':');
    var hour = parseInt(time[0]);
    var minute = parseInt(time[1]);
    if (isBarTime) {
        if (hour >= 6 && hour < 22) {
            hour = 6;
            minute = 0;
        }
    } else {
        if (hour < 12 || hour >= 24) {
            hour = 24;
            minute = 0;
        }
    }
    var h_pos = baseImageInfo.timePosTable.get(hour);
    var h_pos_next = baseImageInfo.timePosTable.get(hour + 1);
    if (!h_pos_next) return h_pos; //undefined
    return h_pos + minute * (h_pos_next - h_pos) / 60;
}
