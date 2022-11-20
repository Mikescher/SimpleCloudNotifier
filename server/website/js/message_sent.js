
window.addEventListener("load", function ()
{
	const qp = new URLSearchParams(window.location.search);

	const spanQuota1 = document.getElementById("insQuota1");
	const spanQuota2 = document.getElementById("insQuota2");
	const linkSucc   = document.getElementById("succ_link");
	const linkErr    = document.getElementById("err_link");

	spanQuota1.innerText = qp.get('quota_remain') ?? 'ERR';
	spanQuota2.innerText = qp.get('quota_max') ?? 'ERR';

	const preset_user_id  = qp.get('preset_user_id')  ?? 'ERR';
	const preset_user_key = qp.get('preset_user_key') ?? 'ERR';

	linkSucc.setAttribute("href", "/?preset_user_id="+preset_user_id+"&preset_user_key="+preset_user_key);

	if (qp.get("ok") === "1") {

		linkSucc.classList.remove('display_none');
		linkErr.classList.add('display_none');

	} else {

		linkSucc.classList.add('display_none');
		linkErr.classList.remove('display_none');

	}

}, false);