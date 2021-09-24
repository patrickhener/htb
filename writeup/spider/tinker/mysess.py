#!/usr/bin/env python3
import flask_unsign

def getCookie(payload):
	session_cookie = flask_unsign.sign(
		value={"cart_items":[],"uuid":payload},
		legacy=True,
		secret="Sup3rUnpredictableK3yPleas3Leav3mdanfe12332942"
	)

	return session_cookie