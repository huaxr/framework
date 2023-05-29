#!/bin/bash
#不变更supervisor配置：sh /home/www/talcamp/deploy.sh talcamp talcamp 0
#变更supervisor配置：sh /home/www/talcamp/deploy.sh talcamp talcamp 1
# worker的命名为worker_talcamp

if [ x"$1" = x ]; then
    echo "projectname param err"
    exit 1
fi

if [ x"$2" = x ]; then
    echo "service param err"
    exit 1
fi

if [ x"$4" = x ]; then
    echo "not found env param"
else
    env_file="/home/www/$1/conf/conf_$4.ini"
    if [ ! -f "$env_file" ]; then
        echo "$env_file config not found"
        exit 1
    fi
    conf_file="/home/www/$1/conf/conf.ini"
    cp -f ${env_file} ${conf_file}
fi
##同时启动worker和service 多个主机名用空格分割
worker_service_hostname=''
##只启动worker 多个主机名用空格分割
worker_hostname=''

projectname=$1
servicename=$2
iscover=$3
function runserver(){
    project=$1
    service=$2
    flag=$3
    supervisorini="/etc/supervisor/${service}.ini"
    projectini="/home/www/${project}/conf/${service}.ini"

    if [ ! -f "$projectini" ]; then
        echo "$projectini config not found"
        return
    fi
    if [ ! -d "/home/logs/xeslog/${service}" ]; then
        mkdir -p /home/logs/xeslog/${service}
    fi

    if [ ! -f "$supervisorini" ]; then ##如果文件不存在 复制文件
        cp -f ${projectini} ${supervisorini}
        supervisorctl update
        return
    elif [ "$flag" = "1" ]; then  ##如果配置不一样变更配置
        checksum=`md5sum "${supervisorini}" | cut -d " " -f1`
        checksum1=`md5sum "${projectini}" | cut -d " " -f1`
        echo ${checksum}
        echo ${checksum1}
        if [ "$checksum" = "$checksum1" ]; then
            echo "ini not change"
        else
            echo "ini change"
            cp -f ${projectini} ${supervisorini}
            supervisorctl update

            #配置变更重启服务
            supervisorctl restart ${service}
            return
        fi
    fi
    #默认动作信号处理
    supervisorctl  signal SIGUSR2 ${service}
}
isworker=0
isall=0
##worker机发布
for name in $worker_hostname
do
	if [ $(hostname) = $name ];then
	    isworker=1
	fi
done
##测试机发布
for name in $worker_service_hostname
do
    if [ $(hostname) = $name ];then
        isall=1
    fi
done

if [ $isall = 1 ]; then
    runserver ${projectname} "worker_${servicename}" ${iscover}
    runserver ${projectname} ${servicename} ${iscover}
elif [ $isworker = 1 ];then
    runserver ${projectname} "worker_${servicename}" ${iscover}
else
    runserver ${projectname} ${servicename} ${iscover}
fi
