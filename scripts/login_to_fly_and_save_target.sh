#!/bin/sh

fly login \
	--target='tpe-cf-mgmt' \
	--concourse-url='https://tpe-concourse-rock.eng.vmware.com/' \
	--team-name='identity-and-credentials'

fly -t tpe-cf-mgmt sync
echo "Fly version:" $(fly -v)