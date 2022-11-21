
function send()
{
	let me = document.getElementById("btnSend");
	if (me.classList.contains("btn-disabled")) return;

	me.innerHTML = "<div class=\"spinnerbox\"><div class=\"spinner primary\"></div></div>";

	me.classList.add("btn-disabled");

	let uid = document.getElementById("uid");
	let key = document.getElementById("ukey");
	let tit = document.getElementById("tit");
	let cnt = document.getElementById("cnt");
	let pio = document.getElementById("prio");
	let cha = document.getElementById("chan");

	uid.classList.remove('input-invalid');
	key.classList.remove('input-invalid');
	msg.classList.remove('input-invalid');
	cnt.classList.remove('input-invalid');
	pio.classList.remove('input-invalid');

	let data = new FormData();
	data.append('user_id', uid.value);
	data.append('user_key', key.value);
	if (tit.value !== '') data.append('title', tit.value);
	if (cnt.value !== '') data.append('content', cnt.value);
	if (pio.value !== '') data.append('priority', pio.value);
	if (cha.value !== '') data.append('channel', cha.value);

	let xhr = new XMLHttpRequest();
	xhr.open('POST', '/send.php', true);
	xhr.onreadystatechange = function ()
	{
		if (xhr.readyState !== 4) return;

		console.log('Status: ' + xhr.status);
		if (xhr.status === 200 || xhr.status === 401 || xhr.status === 403 || xhr.status === 412)
		{
			let resp = JSON.parse(xhr.responseText);
			if (!resp.success || xhr.status !== 200)
			{
				if (resp.errhighlight === 101) uid.classList.add('input-invalid');
				if (resp.errhighlight === 102) key.classList.add('input-invalid');
				if (resp.errhighlight === 103) tit.classList.add('input-invalid');
				if (resp.errhighlight === 104) cnt.classList.add('input-invalid');
				if (resp.errhighlight === 105) pio.classList.add('input-invalid');
				if (resp.errhighlight === 106) cha.classList.add('input-invalid');

				Toastify({
					text: resp.message,
					gravity: "top",
					positionLeft: false,
					backgroundColor: "#D32F2F",
				}).showToast();
			}
			else
			{
				window.location.href =
					'/message_sent' +
					'?ok=' + 1 +
					'&message_count=' + resp.messagecount +
					'&quota=' + resp.quota +
					'&quota_remain=' + (resp.quota_max-resp.quota) +
					'&quota_max=' + resp.quota_max +
					'&preset_user_id=' + uid.value +
					'&preset_user_key=' + key.value +
					'&preset_channel=' + cha.value;
			}
		}
		else
		{
			Toastify({
				text: 'Request failed: Statuscode=' + xhr.status,
				gravity: "top",
				positionLeft: false,
				backgroundColor: "#D32F2F",
			}).showToast();
		}

		me.classList.remove("btn-disabled");
		me.innerHTML = "Send";
	};
	xhr.send(data);
}

window.addEventListener("load", function ()
{
	const qp = new URLSearchParams(window.location.search);

	let btn = document.getElementById("btnSend");
	let uid = document.getElementById("uid");
	let key = document.getElementById("ukey");
	let tit = document.getElementById("tit");
	let cnt = document.getElementById("cnt");
	let pio = document.getElementById("prio");
	let cha = document.getElementById("chan");

	btn.onclick = function () { send(); return false; };

	if (qp.has('preset_priority')) pio.selectedIndex = parseInt(qp.get("preset_priority"));
	if (qp.has('preset_user_key')) key.value         =          qp.get("preset_user_key");
	if (qp.has('preset_user_id'))  uid.value         =          qp.get("preset_user_id");
	if (qp.has('preset_title'))    tit.value         =          qp.get("preset_title");
	if (qp.has('preset_content'))  cnt.value         =          qp.get("preset_content");
	if (qp.has('preset_channel'))  cha.value         =          qp.get("preset_channel");

}, false);