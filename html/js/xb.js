function LeonInit() {
	$('#editEvent').on('show.bs.modal', function (event) {
		var button = $(event.relatedTarget);
		var eventId = button.data('event-id');
		var modal = $(this);
		//ajax here
		if (0 != eventId) {
			$.ajax({
				url:'/xiaobai/msg/event-info',
				contentType:'application/json',
				data:JSON.stringify({Id:eventId}),
				type:'post',
				dataType:'json',
				async: false,
				success:function(data) {
					modal.find('#event-time').val(data.Time);
					modal.find('#event-address').val(data.Address);
					modal.find('#event-title').val(data.Title);
					modal.find('#event-image img').attr('src', data.Image + "?imageView2/4/w/300/h/300");
					modal.find('#event-desc').val(data.Desc);
					haveChgd = new Array();
				}
			});
		} else {
			modal.find('#event-time').val("");
			modal.find('#event-address').val("");
			modal.find('#event-title').val("");
			modal.find('#event-image img').attr('src', '');
			modal.find('#event-desc').val("");
			haveChgd = new Array();
		}
		// if (modal.find('.modal-header #xb-event-id').text() != eventId) {
		// 	modal.find('.modal-header #xb-event-id').text(eventId);
		// 	modal.find('#chg-img').text(0);
		// }
		modal.find('.modal-header #xb-event-id').text(eventId);
		modal.find('.modal-footer .btn-primary').attr('onclick', "putSave()");
	});

	$('#bigImg').on('show.bs.modal', function (event) {
		var button = $(event.relatedTarget);
		var imgSrc = button.data('img-src');
		var modal = $(this);
		modal.find('img').attr('src', imgSrc);
	});

	uploader = Qiniu.uploader({
		runtimes: 'html5',
		browse_button: 'event-image-add',
		// uptoken_url: '/msg/qiniu-token',  // Ajax请求uptoken的Url，格式过于死板
		uptoken_func: function() {
			// var token = "";
			$.ajax({
				url:'/xiaobai/msg/qiniu-token',
				contentType: 'application/json',
				data:JSON.stringify({Bucket:'xiaobai'}),
				type:'post',
				dataType:'json',
				async:false,
				success:function(data) {
					token = data.Token;
				}
			})
			return token;
		},
		unique_names: true,
		get_new_uptoken: false,
		domain: 'xiaobai',
		container: 'event-image',
		max_file_size: '100mb',
		max_retries: 3,
		chunk_size: '4mb',
		multi_selection: false,
		filters: {
			mime_types: [
				{title:"图片文件", extensions:"jpg,png"}
			]
		},
		auto_start: false,
		init: {
			'FilesAdded': function(up, files) {
				plupload.each(files, function(file) {
					// 文件添加进队列后，处理相关的事情
					console.log("add file:", file.name);
					var preloader = new o.Image();//mOxie
					preloader.onload = function() {
						preloader.downsize(300, 300); //压缩下显示 不影响上传
						var imgsrc = preloader.type=='image/jpeg' ? preloader.getAsDataURL('image/jpeg',80) : preloader.getAsDataURL();
						$('#editEvent #event-image img').attr('src', imgsrc);
						preloader.destroy();
						preloader = null;
					};
					preloader.load(file.getSource());
					// $('#editEvent #event-image #chg-img').text(1);
					if (-1 === haveChgd.indexOf('event-image'))
						haveChgd.push('event-image');
				});
			},
			'FileUploaded': function(up, file, info) {
				var domain = up.getOption('domain');
				var res = JSON.parse(info.response);
				var sourceLink = domain +"/"+ res.key; //上传成功后 域+名
				$('#editEvent #event-image img').attr('src', sourceLink);
				saveEvent($('#editEvent .modal-header #xb-event-id').text());
			},
			'Error': function(up, err, errTip) {
				//上传出错时，处理相关的事情
			}
		}
	});

	function bindChg() {
		var chgId = $(this).attr('id');
		if (-1 === haveChgd.indexOf(chgId))
			haveChgd.push(chgId);
	}
	$('#editEvent #event-time').bind('input propertychange', bindChg);
	$('#editEvent #event-address').bind('input propertychange', bindChg);
	$('#editEvent #event-title').bind('input propertychange', bindChg);
	$('#editEvent #event-desc').bind('input propertychange', bindChg);

	laydate.render({
		elem: '#editEvent #event-time',
		done: function(value, date, endDate) {
			if (-1 === haveChgd.indexOf('event-time'))
				haveChgd.push('event-time');
		}
	});

	if (isWx()) {
		$.ajax({
			url: '/xiaobai/msg/wx-config',
			contentType: 'application/json',
			data: JSON.stringify({Url:encodeURIComponent(location.href.split('#')[0])}),
			type: 'post',
			dataType: 'json',
			success:function(data) {
				var cfg = data;
				cfg.debug = false;
				cfg.jsApiList = ['chooseImage', 'previewImage', 'uploadImage', 'downloadImage'];
				wx.config(cfg);
			}
		});
	}
}
function isWx() {
	if ('1' === $('#isWx').text()) {
		return true;
	} else {
		return false;
	}
}
function putSave() {
	// if (1 == $('#editEvent #event-image #chg-img').text()) {
	if (-1 === haveChgd.indexOf('event-image')) {
		saveEvent($('#editEvent .modal-header #xb-event-id').text());
	} else {
		uploader.start();
	}
}
var haveChgd = new Array();
function saveEvent(eventId) {
	var data = {Id:eventId}//new Object();
	var dom = haveChgd.pop();
	while (undefined !== dom) {
		switch (dom) {
			case 'event-time': {
				data.Time = $('#editEvent #event-time').val();
				break;
			}
			case 'event-address': {
				data.Address = $('#editEvent #event-address').val();
				break;
			}
			case 'event-title': {
				data.Title = $('#editEvent #event-title').val();
				break;
			}
			case 'event-image': {
				data.Image = $('#editEvent #event-image img').attr('src');
				break;
			}
			case 'event-desc': {
				data.Desc = $('#editEvent #event-desc').val();
				break;
			}
			default: {
				break;
			}
		}
		dom = haveChgd.pop();
	}
	// console.log(data);
	$.ajax({
		url: '/xiaobai/msg/event-save',
		contentType: 'application/json',
		data: JSON.stringify(data),
		type: 'post',
		dataType: 'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
		}
	});
}
$(function() {
	LeonInit();

	if (window.matchMedia("(max-width: 700px)").matches) {
		$('.div-img').each(mobileFold);
	} else {
		// 随机变换图片
		$('.div-img').each(randomChange);
	}
})
function mobileFold() {
	var i = parseInt($(this).attr('id'), 10);
	// $(this).css('top', 300+i*70+'px');
	$(this).children('img,p,div').addClass('hidden');
	$(this).attr('onclick', 'showActive('+i+')');
}
function showActive(index) {
	$('.div-img#'+index).children('img,p,div').toggleClass('hidden').toggleClass('show');
}
function randomChange() {
	var rR = parseInt(11*Math.random()-5, 10);
	var rTx = parseInt(21*Math.random()-10, 10);
	var rTy = parseInt(21*Math.random()-10, 10);
	var rS = parseInt(5*Math.random()+8, 10)/10;
	$(this).css('transform', 'rotate('+rR+'deg) translate('+rTx+'px, '+rTy+'px) scale('+rS+', '+rS+')');
}
function toBig(index) {
	$('.div-img-box').removeClass('show').addClass('hidden');
	$('.big-img-box').removeClass('hidden').addClass('show');
	var swiper = new Swiper('.swiper-container', {
		initialSlide: index,
		zoom: true,
		nextButton: '.swiper-button-next',
		prevButton: '.swiper-button-prev'
	});
}
function toSmall() {
	$('.big-img-box').removeClass('show').addClass('hidden');
	$('.div-img-box').removeClass('hidden').addClass('show');
}
