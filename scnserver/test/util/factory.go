package util

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"gopkg.in/loremipsum.v1"
	"testing"
	"time"
)

// # Generated by https://chat.openai.com/chat
// ===========================================
//
//	     Create me a list of 32 example notification messages that I can use in unit tests.
//	     Every notification message contains a title and a content.
//	     Do not to repeat the same words in every message.
//	     Vary the length of the content from short sentences to multiple sentences.
//
//       Create me a list of 8 creative and realistic usernames.
//
//       Create me a list of 8 imaginary phone models
//
//       Create me a list of 8 names for your phone
//
//       Create me a list of 8 discord channel names
//

type msgex struct {
	User       int
	Channel    string
	SenderName string
	Priority   int
	Key        int
	Title      string
	Content    string
	TSOffset   time.Duration
}

type userex struct {
	Idx          int64
	WithClient   bool
	Username     string
	AgentModel   string
	AgentVersion string
	ClientType   string
	FCMTok       string
	ProTok       string
}

type clientex struct {
	User         int
	AgentModel   string
	AgentVersion string
	ClientType   string
	FCMTok       string
}

type Userdat struct {
	UID      string
	SendKey  string
	AdminKey string
	ReadKey  string
}

const PX = -1
const P0 = 0
const P1 = 1
const P2 = 2

const AKEY = 0
const SKEY = 1

var userExamples = []userex{
	{0, true, "", "Starfire", "2.0", "IOS", "FCM_TOK_EX_001", ""},
	{1, true, "", "Galaxy Quest", "2022", "ANDROID", "FCM_TOK_EX_002", ""},
	{2, true, "Dreamer23", "Ocean Explorer", "737edc01", "IOS", "FCM_TOK_EX_003", ""},
	{3, true, "CreativeGenius", "Snow Leopard", "1.0.1.99~3", "ANDROID", "FCM_TOK_EX_004", "ANDROID|v2|PURCHASED:PRO_TOK_001"},
	{4, true, "WanderingSoul", "Ocean Explorer", "737edc01", "IOS", "FCM_TOK_EX_005", ""},
	{5, true, "", "Ocean Explorer", "737edc01", "IOS", "FCM_TOK_EX_006", ""},
	{6, true, "BoldExplorer", "Cyber Nova", "Cyber 4", "IOS", "FCM_TOK_EX_007", ""},
	{7, true, "ImaginationKing", "Galaxy Quest", "2023.1", "ANDROID", "FCM_TOK_EX_008", ""},
	{8, true, "", "Galaxy Quest", "2023.1", "ANDROID", "FCM_TOK_EX_009", ""},
	{9, true, "UniqueUnicorn", "Galaxy Quest", "2023.1", "ANDROID", "FCM_TOK_EX_010", ""},
	{10, false, "", "", "", "", "", ""},
	{11, false, "", "", "", "", "", "ANDROID|v2|PURCHASED:PRO_TOK_002"},
	{12, true, "NoMessageUser", "Ocean Explorer", "737edc01", "IOS", "FCM_TOK_EX_014", ""},
	{13, false, "EmptyUser", "", "", "", "", ""},
	{14, true, "ChanTester1", "StarfireXX", "1.x", "IOS", "FCM_TOK_EX_012", ""},
	{15, true, "ChanTester2", "StarfireXX", "1.x", "IOS", "FCM_TOK_EX_013", ""},
	{16, true, "PagTester1", "StarfireXX", "1.x", "ANDROID", "FCM_TOK_EX_016", ""},
}

var clientExamples = []clientex{
	{2, "GalaxySurfer", "Triple-XXX", "IOS", "FCM_TOK_EX2_001"},
	{2, "GalaxySurfer", "Triple-XXX", "IOS", "FCM_TOK_EX2_002"},
	{4, "Thunder-Bolt-4$", "#12", "ANDROID", "FCM_TOK_EX_005"},  // overwrites FCM from first client
	{6, "GalaxySurfer", "Triple-XXX", "IOS", "FCM_TOK_EX2_002"}, // overwrites FCM from user 2 - client 3 (second extra)
	{9, "GalaxySurfer", "Triple-XXX", "IOS", "FCM_TOK_EX2_004"},
	{9, "DreamWeaver", "Triple-XXX", "IOS", "FCM_TOK_EX2_005"},
	{9, "Galaxy Quest", "2023.1", "ANDROID", "FCM_TOK_EX2_006"},
	{9, "Galaxy Quest", "2023.2", "ANDROID", "FCM_TOK_EX2_006"}, // overwrites FCM from (previous) client 4 (3rd extra)
}

var messageExamples = []msgex{
	{0, "Chatting Chamber", "Mobile Mate", P1, AKEY, "New message from John Doe", "", 0},
	{0, "Chatting Chamber", "", P2, SKEY, "Upcoming event", "Don't forget to attend the staff meeting tomorrow at 9:00am", timeext.FromHours(-10.28)},
	{0, "Unicôdé Häll \U0001f92a", "Mobile Mate", P0, SKEY, "System update", "We will be performing maintenance on the server tonight at 11:00pm. The system may be unavailable for up to an hour.", 0},
	{0, "", "Pocket Pal", P2, SKEY, "Reminder", "Your payment is due by the end of the week", 0},
	{0, "", "", P0, SKEY, "New feature available", "We've added a new feature that allows you to save your favorite items for later", 0},
	{0, "Promotions", "", PX, AKEY, "Security alert", "We've detected unusual activity on your account. Please reset your password as soon as possible", 0},
	{0, "Reminders", "", P1, AKEY, "Important notice", "The office will be closed on Friday for a company-wide event", timeext.FromHours(2.72)},
	{0, "Reminders", "Mobile Mate", PX, AKEY, "Urgent", "There has been a power outage in the building. Please evacuate the premises immediately", 0},
	{0, "Reminders", "", PX, AKEY, "Weather alert", "A severe storm is expected to hit the area tonight. Please take necessary precautions", 0},
	{0, "", "", P0, SKEY, "Congratulations", "You have been selected as Employee of the Month. Please come to the front desk to pick up your prize", 0},
	{0, "", "", PX, AKEY, "Attention", "The water cooler is empty. Could someone please refill it?", timeext.FromHours(-11.29)},
	{0, "Chatting Chamber", "Mobile Mate", P2, SKEY, "Important", "All employees are required to complete a safety training course by the end of the month", 0},
	{0, "", "", P1, AKEY, "FAQ Update", Lipsum(10001, 1), 0},
	{0, "", "", PX, AKEY, "Notice", "There will be a fire drill at 10:00am tomorrow. Please follow the instructions of the fire marshal", 0},
	{0, "", "Cellular Confidant", P2, SKEY, "Invitation", "You are invited to a celebration in honor of our 10-year anniversary. The party will be held on Friday at 7:00pm", 0},
	{0, "", "", P0, SKEY, "Deadline reminder", "Please remember to submit your project proposal by the end of the day \U0001f638", 0},
	{0, "Reminders", "", PX, AKEY, "Attention - The copier is out of toner", "", 0},
	{0, "Reminders", "Cellular Confidant", P2, SKEY, "Reminder", "Don't forget to clock in before starting your shift", timeext.FromHours(0.40)},
	{0, "Reminders", "Cellular Confidant", P1, AKEY, "Important", "There will be a company-wide meeting on Monday at 9:00am in the conference room", timeext.FromHours(23.15)},
	{0, "", "", P2, SKEY, "System update", "We will be performing maintenance on the server tonight at 11:00pm. The system may be unavailable for up to an hour. Please save any unsaved work before then", 0},
	{0, "Promotions", "Pocket Pal", P0, SKEY, "Attention - The first aid kit is running low on supplies.", "", 0},
	{0, "Promotions", "Pocket Pal", PX, AKEY, "Urgent", "We have received a complaint about a safety hazard in the workplace. Please address the issue immediately", 0},

	{1, "", "", P1, AKEY, "New message from Jane Doe", "Hey, what's up?", 0},
	{1, "", "", P2, SKEY, "Reminder: Meeting at 3 PM", "Don't forget to join the meeting in the conference room.", 0},
	{1, "", "", P0, SKEY, "Urgent: Action required", "Please review and respond to this important email as soon as possible.", timeext.FromHours(-1.62)},
	{1, "", "", P2, SKEY, "Congratulations!", "You have successfully completed the first step in our onboarding process.", 0},
	{1, "", "", P0, SKEY, "Notice: Maintenance scheduled", "The server will be down for maintenance tonight from 10 PM to 2 AM. Please save your work and log off before then.", 0},
	{1, "private", "", PX, AKEY, "Security alert", "We have detected suspicious activity on your account. Please reset your password immediately to protect your information.", 0},
	{1, "", "", P1, AKEY, "New follower on Twitter", "You have a new follower on Twitter! Check out their profile and see if you want to follow them back.", 0},
	{1, "private", "", PX, AKEY, "Weather alert", "A severe storm is expected to hit the area tonight. Please take shelter and stay safe.", timeext.FromHours(-5.18)},
	{1, "private", "", PX, AKEY, "Free trial ending", "Your free trial of our service is ending in three days. Upgrade now to continue using it.", 0},
	{1, "", "", P0, SKEY, "Task assigned", "You have been assigned a new task: complete the report by Friday at 5 PM.", 0},
	{1, "", "", PX, AKEY, "Event reminder", "Don't forget to join us for the company picnic this Saturday from 12 PM to 3 PM at the park.", 0},
	{1, "private", "", P2, SKEY, "Shipping update", "Your order has shipped and is on its way. Track your package and expect delivery within the next three days.", 0},

	{2, "", "", P1, AKEY, "New feature available", "We have added a new feature to our app! Check it out and let us know what you think.", 0},
	{2, "", "", PX, AKEY, "Payment overdue", "Your payment is overdue. Please make the payment as soon as possible to avoid late fees.", 0},
	{2, "Ü", "", P2, SKEY, "Account suspended", "Your account has been suspended for violating our terms of service. Please contact us to resolve this issue.", 0},
	{2, "Ö", "", P0, SKEY, "Survey invitation", "We would like to invite you to participate in a survey about your experience with our product. Your feedback is valuable to us.", timeext.FromHours(4.66)},
	{2, "", "", PX, AKEY, "Contest winner", "Congratulations! You are the winner of our latest contest. Please contact us to claim your prize.", 0},
	{2, "", "", P2, SKEY, "Appointment confirmation", "This is a confirmation of your upcoming appointment on Friday at 9 AM. Please reply to confirm or reschedule.", 0},
	{2, "Ö", "", P1, AKEY, "Referral program", "Refer a friend to our service and earn a $20 credit. Share your referral code with them and have them sign up using it.", 0},
	{2, "", "", P2, SKEY, "Price change", "We have changed our pricing for our product. Check out the updated pricing and let us know if you have any questions.", 0},
	{2, "Ä", "", P0, SKEY, "New blog post", "We have published a new blog post on our website. Check it out to learn more about our latest product release.", 0},
	{2, "Ä", "", PX, AKEY, "Support ticket update", "We have received your support ticket and are working on a solution. Please expect a response within the next 24 hours.", 0},
	{2, "Ä", "", P2, SKEY, "Order confirmation", "Thank you for your order! Your order number is 12345. We will send a confirmation email when your order has shipped.", 0},
	{2, "Ü", "", P0, SKEY, "Feedback request", "We value your opinion. Please take a few minutes to complete our survey and let us know how we can improve our service.", 0},
	{2, "", "", P1, AKEY, "New product launch", "We are excited to announce the launch of our newest product! Check it out and let us know what you think.", 0},
	{2, "", "", P0, SKEY, "Free shipping", "Enjoy free shipping on your", 0},

	{3, "\U0001f5ff", "", PX, AKEY, "New message received", "You have a new message from John Doe.", 0},
	{3, "", "", P1, AKEY, "Meeting reminder", "Don't forget, your meeting with the sales team is at 10 AM tomorrow.", 0},
	{3, "", "", PX, AKEY, "Payment confirmation", "Your payment of $100 has been successfully processed. Thank you for your business.", 0},
	{3, "", "", P2, SKEY, "Task completed", "Your task \"Update website content\" has been completed and is ready for review.", 0},
	{3, "Innovations", "", PX, AKEY, "Invitation to join a group", "You have been invited to join the \"Marketing Team\" group on our collaboration platform.", 0},
	{3, "", "", P2, SKEY, "Password reset", Lipsum(10002, 1), 0},
	{3, "", "", P2, SKEY, "Low battery alert", Lipsum(10003, 2), 0},
	{3, "Innovations", "", P2, SKEY, "System update available", Lipsum(10004, 5), 0},
	{3, "", "", P2, SKEY, "Appointment confirmation", "Your appointment for a physical exam on Monday, March 15th at 10 AM has been confirmed.", 0},
	{3, "\U0001f5ff", "", P2, SKEY, "Order shipped", "Your order #123456 has been shipped and is on its way to your address.", 0},
	{3, "", "", P2, SKEY, "Order cancelled", "Your order #123456 has been cancelled. We apologize for any inconvenience this may have caused.", 0},
	{3, "", "", P2, SKEY, "Event reminder", "Don't forget, the company holiday party is tomorrow at 6 PM. We hope to see you there!", 0},
	{3, "Reminders", "", PX, AKEY, "Account verification", "", timeext.FromHours(1.15)},
	{3, "Reminders", "", PX, AKEY, "Overdue payment", "", 0},
	{3, "Reminders", "", P2, SKEY, "Security alert", "We have detected suspicious activity on your account. Please take the necessary steps to secure your account.", timeext.FromHours(0.80)},
	{3, "Reminders", "", PX, AKEY, "Product back in stock", Lipsum(10001, 6), 0},
	{3, "", "", PX, AKEY, "Connection lost", "Your device has lost its connection to the internet. Please check your network settings and try again.", 0},
	{3, "", "", P2, SKEY, "Subscription renewal", "Your subscription is set to renew in one week. Please update your payment information to avoid any interruption in service.", 0},
	{3, "", "", PX, AKEY, "Work order assigned", "You have been assigned a new work order #123456. Please review the details and complete the task as soon as possible.", 0},
	{3, "Innovations", "", P2, SKEY, "Scheduled maintenance", "", 0},
	{3, "Innovations", "", P2, SKEY, "Payment declined", "Your payment for invoice #123456 has been declined. Please update your payment information and try again.", 0},
	{3, "Innovations", "", P1, AKEY, "New follower", "You have a new follower on our platform. Welcome them with a message and start building your network.", 0},
	{3, "", "", P1, AKEY, "Account suspended", "", 0},
	{3, "\U0001f5ff", "", P0, SKEY, "Request for feedback", "We value your feedback and would love to hear your thoughts on your recent experience with our platform.", 0},
	{3, "\U0001f5ff", "", P0, SKEY, "Task deadline approaching", "Your task \"Write blog post\" is due in three days. Please make sure to complete it on time", 0},

	{4, "", "", P0, SKEY, "Server maintenance", "The server will be offline for maintenance on Tuesday, January 5th at 10pm EST", 0},
	{4, "", "", P0, SKEY, "New feature update", "A new feature has been added to the server. Please check the changelog for details", 0},
	{4, "", "Server0", PX, AKEY, "Security alert", "There has been a security breach on the server. Please change your password immediately", timeext.FromHours(-4.90)},
	{4, "", "Server0", P0, SKEY, "Server upgrade", "The server has been upgraded with improved performance and security features", timeext.FromHours(6.29)},
	{4, "", "", P2, SKEY, "Scheduled downtime", "The server will be offline for scheduled downtime on Friday, January 8th at 8am EST", 0},
	{4, "", "Server0", P0, SKEY, "Server outage", "There is currently a server outage. We are working to resolve the issue as soon as possible", 0},
	{4, "", "", P0, SKEY, "User account update", "Your user account has been updated with new features and improved security measures", 0},
	{4, "", "", P1, AKEY, "Server status update", "The server is currently experiencing higher than normal traffic. We apologize for any inconvenience", 0},
	{4, "", "Server0", P0, SKEY, "Server upgrade", "The server has been upgraded again", timeext.FromHours(6.19)},

	{5, "", "localhost", P1, AKEY, "New server release", "A new version of the server has been released. Please update to the latest version to ensure optimal performance", 0},
	{5, "Test1", "localhost", P1, AKEY, "Server maintenance schedule", "The server will be undergoing regular maintenance every Tuesday at 10pm EST", timeext.FromHours(12.45)},
	{5, "Test2", "example.com", P0, SKEY, "Server outage notification", "We apologize for the inconvenience, but the server is currently experiencing an outage", 0},
	{5, "Test3", "example.com", PX, AKEY, "Server performance update", "The server is currently experiencing improved performance thanks to recent upgrades", timeext.FromHours(-0.18)},
	{5, "Test4", "example.org", P1, AKEY, "Server security patch", "A security patch has been applied to the server to improve its security measures", 0},
	{5, "Test5", "example.org", PX, AKEY, "Server downtime schedule", "The server will be offline for scheduled downtime on the first Friday of every month at 8am EST", 0},

	{6, "", "server2", P0, SKEY, "Server outage resolution", "The server outage has been resolved and the server is now back online", 0},
	{6, "", "server1", P0, SKEY, "New server features", "The server has been updated with new features and improved functionality", 0},
	{6, "", "server2", P2, SKEY, "Server traffic update", "The server is currently experiencing high traffic levels. We apologize for any delays", 0},
	{6, "Security", "server2", P0, SKEY, "Server security breach", "There has been a security breach on the server. Please change your password and update your security settings", 0},
	{6, "", "server1", P1, AKEY, "Server maintenance notification", "The server will be offline for maintenance on Wednesday, January 6th at 10pm EST", 0},
	{6, "", "server1", P0, SKEY, "Server outage warning", "We are experiencing server issues and may need to take the server offline for maintenance", timeext.FromHours(6.65)},
	{6, "", "server1", P2, SKEY, "Server performance improvement", "Thanks to recent upgrades, the server is now performing better than ever", 0},
	{6, "", "server1", PX, AKEY, "Server security update", "The server has been updated with the latest security patches and enhancements", 0},
	{6, "", "server1", P1, AKEY, "Server downtime schedule change", "The server downtime schedule has been changed to every other Friday at 8am EST", 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20001, 1), 0},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20002, 1), 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20003, 1), 0},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20004, 1), 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20005, 1), 0},
	{6, "Lipsum", "", P1, AKEY, "Lorem Ipsum", Lipsum(20006, 1), 0},
	{6, "Lipsum", "", P1, AKEY, "Lorem Ipsum", Lipsum(20007, 1), timeext.FromHours(-3.39)},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20008, 1), 0},
	{6, "Lipsum", "", PX, AKEY, "Lorem Ipsum", Lipsum(20009, 1), 0},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20010, 1), 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20011, 1), 0},
	{6, "Lipsum", "", PX, AKEY, "Lorem Ipsum", Lipsum(20012, 1), 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20013, 1), 0},
	{6, "Lipsum", "", P2, SKEY, "Lorem Ipsum", Lipsum(20014, 1), timeext.FromHours(-2.33)},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20015, 1), 0},
	{6, "Lipsum", "", P0, SKEY, "Lorem Ipsum", Lipsum(20016, 1), 0},

	{7, "", "localhost", P2, SKEY, "Server outage resolution update", "We are still working on resolving the server outage and will provide updates as soon as possible", 0},
	{7, "", "localhost", P0, SKEY, "New server release update", "A new update for the server has been released. Please update to the latest version for optimal performance", 0},
	{7, "", "localhost", P2, SKEY, "Server traffic warning", "The server is experiencing high traffic levels and may be slow. We apologize for any inconvenience", 0},
	{7, "", "localhost", P1, AKEY, "Server security alert", "There has been a potential security breach on the server. Please update your password and security settings immediately", 0},
	{7, "", "localhost", PX, AKEY, "Server maintenance reminder", "Don't forget, the server will be offline for maintenance on Thursday, January 7th at 10pm EST", 0},
	{7, "", "localhost", P1, AKEY, "Server outage status", "The server outage is ongoing and we are working to resolve the...", 0},

	{8, "", "", PX, AKEY, "Get your free trial now!", "Sign up for our exclusive offer and get access to our premium features for a limited time!", 0},
	{8, "", "", PX, AKEY, "Limited time offer", "Hurry and take advantage of our discounted rates before they expire!", 0},
	{8, "", "", PX, AKEY, "Unbeatable deals", "Get the best prices on the hottest products today! Act fast, these deals won't last long.", 0},
	{8, "", "", P2, SKEY, "One-click signup", "Join our mailing list and get instant access to exclusive offers and discounts.", 0},
	{8, "", "", P0, SKEY, "Don't miss out", "Join now and get access to our members-only perks and benefits.", 0},
	{8, "", "", P2, SKEY, "Sign up and save", "Get instant savings when you join our email list and be the first to hear about our special deals and promotions.", 0},
	{8, "", "", P0, SKEY, "Exclusive offer", "Sign up now and get a free gift with your first purchase!", timeext.FromHours(10.81)},

	{9, "", "Max", P0, SKEY, "Special discount", "", 0},
	{9, "", "Tim", P1, AKEY, "Huge savings", "", 0},
	{9, "", "Vincent", P0, SKEY, "Insider access", "", 0},

	{10, "", "", PX, AKEY, "Join the club", "Become a member today and get access to exclusive perks and benefits, plus get a free gift with your first purchase.", 0},
	{10, "", "", P2, SKEY, "Join now and save", "Sign up for our email list and get instant access to exclusive offers and discounts on your favorite products.", 0},
	{10, "", "", P1, AKEY, "Hurry, limited time offer", "Sign up now and get a free trial of our premium services before the offer expires!", 0},
	{10, "", "", P2, SKEY, "Limited time only", "Sign up now and get a free gift with your first purchase, plus get access to exclusive deals and discounts.", 0},
	{10, "", "", PX, AKEY, "Sign up and save big", "Join our email list and get instant access to special offers and discounts on top brands and products.", 0},
	{10, "", "", P0, SKEY, "Exclusive membership", "Join now and get access to our members-only perks and benefits, plus get a free gift with your first purchase.", 0},
	{10, "", "", P2, SKEY, "Don't miss out on savings", "Sign up for our email list and get instant access to exclusive offers and discounts on your favorite products.", 0},
	{10, "", "", P1, AKEY, "Sign up and get a free gift", "Join now and get a free gift with your first purchase, plus get access to exclusive deals and discounts.", 0},
	{10, "", "", P2, SKEY, "Limited time offer", "Sign up now and get a free trial of our premium services before the offer expires!", 0},
	{10, "", "", P0, SKEY, "Join now and save", "Sign up for our email list and get instant access to exclusive offers and discounts on top brands and products.", 0},
	{10, "", "", PX, AKEY, "Exclusive offer", "Join now and get access to our members-only perks and benefits, plus get a free gift with your first purchase.", 0},

	{11, "Promotions", "localhost", P2, SKEY, "New Product Launch: Introducing Our Latest Innovation", "We are excited to announce the release of our newest product, designed to revolutionize the industry. Don't miss out on this game-changing technology.", timeext.FromHours(-12.21)},
	{11, "Promotions", "#S0", P0, SKEY, "Limited Time Offer: Get 50% Off Your Next Purchase", "For a limited time, take advantage of our special offer and get half off your next purchase. Don't miss out on this amazing deal.", 0},
	{11, "Promotions", "#S0", P2, SKEY, "Customer Appreciation Sale: Save Up to 75% on Your Favorite Products", "", 0},
	{11, " Promotions", "", P0, SKEY, "Sign Up for Our Newsletter and Save 10% on Your Next Order", "", 0},
	{11, "Promotions ", "", PX, AKEY, "New Arrivals: Check Out Our Latest Collection", "We've just added new items to our collection and we think you'll love them. Take a look and see what's new in fashion, home decor, and more.", 0},
	{11, "Promotions", "", PX, AKEY, "Join Our Rewards Program and Earn Points on Every Purchase", "Sign up for our rewards program and earn points on every purchase you make. Redeem your points for discounts, free products, and more.", 0},
	{11, "Promotions", "#S0", P0, SKEY, "Seasonal Special: Save on Your Favorite Fall Products", "As the leaves change color and the air gets cooler, we have the perfect products to help you enjoy the season. Take advantage of our special offers and save on your favorite fall products.", 0},
	{11, "Promotions", "192.168.0.1", P1, AKEY, "Refer a Friend and Save on Your Next Order", "Share the love and refer a friend to our store. When they make a purchase, you'll receive a discount on your next order. It's a win-win for both of you.", 0},
	{11, "Promotions", "", P2, SKEY, "Free Shipping on All Orders Over $50", "", 0},
	{11, "Promotions", "", PX, AKEY, "Buy One, Get One 50% Off: Mix and Match Your Favorite Products", "", 0},
	{11, "Promotions", "192.168.0.1", P1, AKEY, "New Customer Coupon: Save $10 on Your First Order", "Welcome to our store! As a new customer, we want to offer you a special discount on your first order. Use the coupon code NEW10 at checkout and save $10 on your purchase.", 0},
	{11, "Promotions", "192.168.0.1", P1, AKEY, "Announcing Our Annual Black Friday Sale", "Mark your calendars and get ready for the biggest sale of the year. Our annual Black Friday sale is coming soon and you won't want to miss out on the amazing deals and discounts.", 0},
	{11, "Promotions", "", PX, AKEY, "Join Our VIP Club and Enjoy Exclusive Benefits", "Sign up for our VIP club and enjoy exclusive benefits like early access to sales, special offers, and personalized service. Don't miss out on this exclusive opportunity.", timeext.FromHours(2.32)},
	{11, "Promotions", "", P2, SKEY, "Summer Clearance: Save Up to 75% on Your Favorite Products", "It's time for our annual summer clearance sale! Save up to 75% on your favorite products, from clothing and accessories to home decor and more.", timeext.FromHours(1.87)},

	{14, "", "", P0, SKEY, "New Product Launch", "We are excited to announce the launch of our new product, the XYZ widget", 0},
	{14, "chan_self_subscribed", "", P0, SKEY, "Important Update", "We have released a critical update", 0},
	{14, "chan_self_unsub", "", P0, SKEY, "Reminder: Upcoming Maintenance", "", 0},

	{15, "", "", P0, SKEY, "New Feature Available", "ability to schedule appointments", 0},
	{15, "chan_other_nosub", "", P0, SKEY, "Account Suspended", "Please contact us", 0},
	{15, "chan_other_request", "", P0, SKEY, "Invitation to Beta Test", "", 0},
	{15, "chan_other_accepted", "", P0, SKEY, "New Blog Post", "Congratulations on your promotion! We are proud", 0},

	{16, "Chan1", "", P2, SKEY, "Lorem Ipsum 01", Lipsum(30001, 1), 0},
	{16, "Chan2", "", P0, SKEY, "Lorem Ipsum 02", Lipsum(30002, 1), 0},
	{16, "Chan1", "", P2, SKEY, "Lorem Ipsum 03", Lipsum(30003, 1), 0},
	{16, "Chan1", "", P0, SKEY, "Lorem Ipsum 04", Lipsum(30004, 1), 0},
	{16, "Chan1", "", P2, SKEY, "Lorem Ipsum 05", Lipsum(30005, 1), 0},
	{16, "Chan1", "", P1, AKEY, "Lorem Ipsum 06", Lipsum(30006, 1), 0},
	{16, "Chan2", "", P1, AKEY, "Lorem Ipsum 07", Lipsum(30007, 1), 0},
	{16, "Chan2", "", P0, SKEY, "Lorem Ipsum 08", Lipsum(30008, 1), 0},
	{16, "Chan1", "", PX, AKEY, "Lorem Ipsum 09", Lipsum(30009, 1), 0},
	{16, "Chan1", "", P0, SKEY, "Lorem Ipsum 10", Lipsum(30010, 1), 0},
	{16, "Chan1", "", P2, SKEY, "Lorem Ipsum 11", Lipsum(30011, 1), 0},
	{16, "Chan2", "", PX, AKEY, "Lorem Ipsum 12", Lipsum(30012, 1), 0},
	{16, "Chan2", "", P2, SKEY, "Lorem Ipsum 13", Lipsum(30013, 1), 0},
	{16, "Chan2", "", P2, SKEY, "Lorem Ipsum 14", Lipsum(30014, 1), 0},
	{16, "Chan2", "", P0, SKEY, "Lorem Ipsum 15", Lipsum(30015, 1), 0},
	{16, "Chan3", "", P0, SKEY, "Lorem Ipsum 16", Lipsum(30016, 1), 0},
	{16, "Chan3", "", P0, SKEY, "Lorem Ipsum 17", Lipsum(30017, 1), 0},
	{16, "Chan3", "", P2, SKEY, "Lorem Ipsum 18", Lipsum(30018, 1), 0},
	{16, "Chan3", "", PX, AKEY, "Lorem Ipsum 19", Lipsum(30019, 1), 0},
	{16, "Chan3", "", P2, SKEY, "Lorem Ipsum 20", Lipsum(30020, 1), 0},
	{16, "Chan2", "", P2, SKEY, "Lorem Ipsum 21", Lipsum(30021, 1), 0},
	{16, "Chan2", "", P0, SKEY, "Lorem Ipsum 22", Lipsum(30022, 1), 0},
	{16, "Chan2", "", P0, SKEY, "Lorem Ipsum 23", Lipsum(30023, 1), 0},
}

type DefData struct {
	User []Userdat
}

func InitDefaultData(t *testing.T, ws *logic.Application) DefData {

	// set logger to buffer, only output if error occured
	success := false
	SetBufLogger()
	defer func() {
		ClearBufLogger(!success)
		if success {
			log.Info().Msgf("Succesfully initialized default data (%d messages, %d users)", len(messageExamples), len(userExamples))
		}
	}()

	baseUrl := "http://127.0.0.1:" + ws.Port

	users := make([]Userdat, 0, len(userExamples))

	// Create Users

	for _, uex := range userExamples {
		body := gin.H{}
		if uex.WithClient {
			body["agent_model"] = uex.AgentModel
			body["agent_version"] = uex.AgentVersion
			body["client_type"] = uex.ClientType
			body["fcm_token"] = uex.FCMTok
		} else {
			body["no_client"] = true
		}
		if uex.Username != "" {
			body["username"] = uex.Username
		}
		if uex.ProTok != "" {
			body["pro_token"] = uex.ProTok
		}

		user0 := RequestPost[gin.H](t, baseUrl, "/api/v2/users", body)
		uid0 := user0["user_id"].(string)
		readtok0 := user0["read_key"].(string)
		sendtok0 := user0["send_key"].(string)
		admintok0 := user0["admin_key"].(string)
		AssertMultiNonEmpty(t, "user0", uid0, readtok0, sendtok0, admintok0)

		users = append(users, Userdat{
			UID:      uid0,
			SendKey:  sendtok0,
			AdminKey: admintok0,
			ReadKey:  readtok0,
		})
	}

	// Create Clients

	for _, cex := range clientExamples {
		body := gin.H{}
		body["agent_model"] = cex.AgentModel
		body["agent_version"] = cex.AgentVersion
		body["client_type"] = cex.ClientType
		body["fcm_token"] = cex.FCMTok
		RequestAuthPost[gin.H](t, users[cex.User].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/clients", users[cex.User].UID), body)
	}

	// Create Messages

	for _, mex := range messageExamples {
		body := gin.H{}
		body["title"] = mex.Title
		body["user_id"] = users[mex.User].UID
		switch mex.Key {
		case AKEY:
			body["key"] = users[mex.User].AdminKey
		case SKEY:
			body["key"] = users[mex.User].SendKey
		}
		if mex.Content != "" {
			body["content"] = mex.Content
		}
		if mex.SenderName != "" {
			body["sender_name"] = mex.SenderName
		}
		if mex.Channel != "" {
			body["channel"] = mex.Channel
		}
		if mex.Priority != PX {
			body["priority"] = mex.Priority
		}
		if mex.TSOffset != 0 {
			body["timestamp"] = (time.Now().Add(mex.TSOffset)).Unix()
		}

		RequestPost[gin.H](t, baseUrl, "/", body)
	}

	// create manual channels

	{
		RequestAuthPost[Void](t, users[9].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", users[9].UID), gin.H{"name": "manual@chan"})
	}

	// Sub/Unsub for Users 12+13

	{
		doUnsubscribe(t, baseUrl, users[14], users[14], "chan_self_unsub")
		doSubscribe(t, baseUrl, users[14], users[15], "chan_other_request")
		doSubscribe(t, baseUrl, users[14], users[15], "chan_other_accepted")
		doAcceptSub(t, baseUrl, users[15], users[14], "chan_other_accepted")
	}

	success = true

	return DefData{User: users}
}

type SingleData struct {
	UID      string
	AdminKey string
	SendKey  string
	ReadKey  string
	ClientID string
}

func InitSingleData(t *testing.T, ws *logic.Application) SingleData {

	// set logger to buffer, only output if error occured
	success := false
	SetBufLogger()
	defer func() {
		ClearBufLogger(!success)
		if success {
			log.Info().Msgf("Succesfully initialized default data (%d messages, %d users)", len(messageExamples), len(userExamples))
		}
	}()

	baseUrl := "http://127.0.0.1:" + ws.Port

	type R struct {
		Clients []struct {
			ClientId string `json:"client_id"`
			UserId   string `json:"user_id"`
		} `json:"clients"`
		ReadKey  string `json:"read_key"`
		SendKey  string `json:"send_key"`
		AdminKey string `json:"admin_key"`
		UserId   string `json:"user_id"`
	}

	r0 := RequestPost[R](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	success = true

	return SingleData{
		UID:      r0.UserId,
		AdminKey: r0.AdminKey,
		SendKey:  r0.SendKey,
		ReadKey:  r0.ReadKey,
		ClientID: r0.Clients[0].ClientId,
	}
}

func doSubscribe(t *testing.T, baseUrl string, user Userdat, chanOwner Userdat, chanInternalName string) {

	if user == chanOwner {

		RequestAuthPost[Void](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", user.UID), gin.H{
			"channel_owner_user_id": chanOwner.UID,
			"channel_internal_name": chanInternalName,
		})

	} else {
		type chanlist struct {
			Channels []gin.H `json:"channels"`
		}

		clist := RequestAuthGet[chanlist](t, chanOwner.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=owned", chanOwner.UID))

		var chandat gin.H
		for _, v := range clist.Channels {
			if v["internal_name"].(string) == chanInternalName {
				chandat = v
				break
			}
		}

		RequestAuthPost[Void](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", user.UID, chandat["subscribe_key"].(string)), gin.H{
			"channel_id": chandat["channel_id"].(string),
		})

	}

}

func doUnsubscribe(t *testing.T, baseUrl string, user Userdat, chanOwner Userdat, chanInternalName string) {

	type chanlist struct {
		Subscriptions []gin.H `json:"subscriptions"`
	}

	slist := RequestAuthGet[chanlist](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?direction=outgoing&confirmation=confirmed", user.UID))

	var subdat gin.H
	for _, v := range slist.Subscriptions {
		if v["channel_internal_name"].(string) == chanInternalName && v["channel_owner_user_id"].(string) == chanOwner.UID {
			subdat = v
			break
		}
	}

	RequestAuthDelete[Void](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%v", user.UID, subdat["subscription_id"]), gin.H{})

}

func doAcceptSub(t *testing.T, baseUrl string, user Userdat, subscriber Userdat, chanInternalName string) {

	type chanlist struct {
		Subscriptions []gin.H `json:"subscriptions"`
	}

	slist := RequestAuthGet[chanlist](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?direction=incoming&confirmation=unconfirmed", user.UID))

	var subdat gin.H
	for _, v := range slist.Subscriptions {
		if v["channel_internal_name"].(string) == chanInternalName && v["subscriber_user_id"].(string) == subscriber.UID {
			subdat = v
			break
		}
	}

	RequestAuthPatch[Void](t, user.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%v", user.UID, subdat["subscription_id"]), gin.H{
		"confirmed": true,
	})

}

func LipsumWord(seed int64, wordcount int) string {
	return loremipsum.NewWithSeed(seed).Words(wordcount)
}

func Lipsum(seed int64, paracount int) string {
	return loremipsum.NewWithSeed(seed).Paragraphs(paracount)
}

func ShortLipsum(seed int64, wcount int) string {
	return loremipsum.NewWithSeed(seed).Words(wcount)
}
func ShortLipsum0(wcount int) string {
	return loremipsum.NewWithSeed(0).Words(wcount)
}
