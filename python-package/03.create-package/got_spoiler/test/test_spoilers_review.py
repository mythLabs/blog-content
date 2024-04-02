import sys
sys.path.append(r'./src')

import unittest
from spoilers import review 

class TestReview(unittest.TestCase):

	def test_review(self):
		self.assertEqual(review.season(1),"Epic")

unittest.main()