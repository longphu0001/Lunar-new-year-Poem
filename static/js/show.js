function clickButton_1(){
	var id = Math.floor((Math.random() * 5) + 1);
	location.href='/show?id='+id;
	console.log('id = ',id)
}