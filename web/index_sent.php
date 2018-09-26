<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="/css/mini-default.min.css">
	<!--<link rel="stylesheet" href="/css/mini-nord.min.css">-->
	<!--<link rel="stylesheet" href="/css/mini-dark.min.css">-->
	<link rel="stylesheet" href="/css/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>


	<form id="mainpnl">

        <div class="fullcenterflex">

            <?php if (isset($_GET['ok']) && $_GET['ok'] === "1" ): ?>

                <a class="card success" href="/index.php?preset_user_id=<?php echo isset($_GET['preset_user_id'])?$_GET['preset_user_id']:'ERR';?>&preset_user_key=<?php echo isset($_GET['preset_user_key'])?$_GET['preset_user_key']:'ERR';?>">
                    <div class="section">
                        <h3 class="doc">Message sent</h3>
                        <p class="doc">Message succesfully sent<br>
                            <?php echo isset($_GET['quota_remain'])?$_GET['quota_remain']:'ERR';?>/<?php echo isset($_GET['quota_max'])?$_GET['quota_max']:'ERR';?> remaining</p>
                    </div>
                </a>

			<?php else: ?>

                <a class="card error" href="/index.php">
                    <div class="section">
                        <h3 class="doc">Failure</h3>
                        <p class="doc">Unknown error</p>
                    </div>
                </a>

            <?php endif; ?>

        </div>

        <a href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier" class="button bordered" id="tl_link"><span class="icn-google-play"></span></a>
        <a href="/index.php" class="button bordered" id="tr_link">Send</a>

        <h1>Simple Cloud Notifier</h1>

    </form>

	<div id="copyinfo">
		<a href="https://www.blackforestbytes.com">&#169; blackforestbytes</a>
	</div>

</body>
</html>