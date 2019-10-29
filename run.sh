#!/bin/bash

usage()
{
cat << EOF
usage: $0 options
OPTIONS:
   -h      Show this message
   -I      Internal DNS record Default: kube-dns.kube-system.svc.cluster.local
   -i      Internal DNS IP address Default: 10.43.0.10
   -E      External DNS record Default: a.root-servers.net
   -e      External DNS IP address Default: 198.41.0.4
   -T      DNS Timeout in seconds Default: 10
EOF
}

VERBOSE=
while getopts .ht:I:i:E:e:T:v. OPTION
do
     case $OPTION in
         h)
             usage
             exit 1
             ;;
         I)
             InternalHost=$OPTARG
             ;;
         i)
             InternalIP=$OPTARG
             ;;
         E)
             ExternalHost=$OPTARG
             ;;
         e)
             ExternalIP=$OPTARG
             ;;
         T)
             Timeout=$OPTARG
             ;;
         ?)
             usage
             exit
             ;;
     esac
done

if [[ -z $InternalHost ]]
then
	InternalHost="kube-dns.kube-system.svc.cluster.local"
fi

if [[ -z $InternalIP ]]
then
        $InternalIP="10.43.0.10"
fi

if [[ -z $ExternalHost ]]
then
        ExternalHost="a.root-servers.net"
fi

if [[ -z $InternalHost ]]
then
        ExternalIP="198.41.0.4"
fi

if [[ -z $Timeout ]]
then
        Timeout=10
fi


while true
do

for ip in $(kubectl describe endpoints kube-dns --namespace=kube-system | grep ' Addresses:' | awk '{print $2}' | sed 's/,/\n/g')
do
	echo "Checking $ip"
	InternalOutput=`timeout "$timeout" dig -short @"$ip" "$InternalHost"`
	if $InternalOutput
	then
		##Internal DNS returned successfully and now we need to check the IP
		if [[ "$InternalOutput" == "$InternalIP" ]]
		then
			echo "Internal DNS is OK"
		else
			echo "Internal DNS returned a bad IP"
		fi
	else
		echo "Internal DNS has timed out"
	fi
        ExternalOutput=`timeout "$timeout" dig -short @"$ip" "$ExternalHost"`
        if $ExternalOutput
        then
                ##External DNS returned successfully and now we need to check the IP
                if [[ "$ExternalOutput" == "$ExternalIP" ]]
                then
                        echo "External DNS is OK"
                else
                        echo "External DNS returned a bad IP"
                fi
        else
                echo "External DNS has timed out"
        fi
done

done
