import psycopg2
import unittest
import os
import requests
import json

requestServiceURL = 'http://127.0.0.1:6666/request'
manageServiceURL = 'http://127.0.0.1:6667/'
taskServiceURL = 'http://127.0.0.1:6668/task'

test_req_type = '300'
conn = psycopg2.connect(database="request", user="postgres", password="postgres", host="192.168.56.132", port=5432)


def checkResponseCode(resp):
	respJson = json.loads(resp.text)
	
	if(respJson['StateCode'] == 100):
		return True, respJson['Msg']
	else:
		return False, respJson['Msg']


class TestRequestHandling(unittest.TestCase):
	'''
	def setUp(self):
		
	def tearDown(self):
	'''

	def test_clean(self):
		parms = {
			'type' : test_req_type
		}
		resp = requests.post(manageServiceURL + "clean", data=parms)
		code, msg = checkResponseCode(resp)
		self.assertTrue(code)
		
		cur = conn.cursor()
		cur.execute("select count(*) from pg_class where relname = 'request_" + test_req_type + "'")
		row = cur.fetchone()
		self.assertEqual(row[0], 0)
		
		cur.execute("select count(*) from pg_class where relname = 'req_state_" + test_req_type + "'")
		row = cur.fetchone()
		self.assertEqual(row[0], 0)
		
		
	def test_createRequestTable(self):
		parms = {
			'type' : test_req_type
		}
		
		resp = requests.post(manageServiceURL+"addRequestType", data=parms)
		code, msg = checkResponseCode(resp)
		self.assertTrue(code)
		
		cur = conn.cursor()
		cur.execute("select count(*) from pg_class where relname = 'request_" + test_req_type + "'")
		row = cur.fetchone()
		self.assertEqual(row[0], 1)
		
		cur.execute("select count(*) from pg_class where relname = 'req_state_" + test_req_type + "'")
		row = cur.fetchone()
		self.assertEqual(row[0], 1)
	
	
	def test_insertNewRequest(self):
		sub = True
		noticeAddr = "notice me HERE"
		body = "a new request--aaaa"
		
		parms = {
			'type' : test_req_type,
			'subscribe' : sub,
			'noticeaddr' : noticeAddr,
			'body': body
		}
		
		resp = requests.post(requestServiceURL, data=parms)
		code, msg = checkResponseCode(resp)
		self.assertTrue(code)
		
		#check request table
		cur = conn.cursor()
		cur.execute("SELECT *  from request_" + test_req_type)
		rows_request = cur.fetchall()
		self.assertEqual(len(rows_request), 1)
		
		row = rows_request[0]
		
		self.assertEqual(row[1], sub)
		self.assertEqual(row[2], noticeAddr)
		self.assertEqual(row[3], body)
		
		#check req_state table
		cur.execute("SELECT *  from req_state_" + test_req_type)
		rows_req_state = cur.fetchall()
		self.assertEqual(len(rows_req_state), 1)
		
		
	def test_getTask(self):	
		resp = requests.post(taskServiceURL+"?type=" + test_req_type)
		code, msg = checkResponseCode(resp)
		self.assertTrue(code)	
		print(msg)
		
if __name__ == '__main__':
	suite = unittest.TestSuite()
	tests = [TestRequestHandling("test_clean"),
		TestRequestHandling("test_createRequestTable"),
		TestRequestHandling("test_insertNewRequest"),
		TestRequestHandling("test_getTask")]
	suite.addTests(tests)
	
	runner = unittest.TextTestRunner(verbosity=2)
	runner.run(suite)
	