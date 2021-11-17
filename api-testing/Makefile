setup:
	pip3 install .

test:
	pytest

report:
	pytest --alluredir=allure-reports

view: report
	allure serve allure-reports