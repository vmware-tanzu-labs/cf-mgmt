#!/bin/sh

fly login \
	--target='tpe-cf-mgmt' \
	--concourse-url='https://tpe-concourse-rock.acc.broadcom.net/' \
	--team-name='identity-and-credentials'

fly -t tpe-cf-mgmt sync
echo "Fly version:" $(fly -v)