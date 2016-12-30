#!/bin/bash

OUTFILE=/tmp/dump.csv
GPLTFile=/tmp/plot.plt
USER=$1
PASSWORD=$2
DATABASE=$3

# $1 = data output file.
# $2 = # days to pull
# $3 = Username
# $4 = password
# $5 = database
function extract() {
	echo using sudo to remove old export 
	sudo rm -fr $1

	#Check out http://mamchenkov.net/wordpress/2011/04/27/mysql-export-csv-into-outfile-triggers-access-denied-error/
	# This requires a global USE permission in mysql;:  EG
	# USE database;
	# UPDATE user SET File_priv = 'Y' WHERE User = 'your-mysql-user';
	# FLUSH PRIVILEGES;
	echo Extracting the last days worth of data
	mysql -u $3 -h localhost -p$4 -D $5 -e "select created,battery,humidity,ihumidity,pressure,temperature,temperatureExt INTO OUTFILE '$1' FIELDS TERMINATED BY ',' LINES TERMINATED BY '\\n' from wxstation WHERE created > DATE_SUB(CURRENT_TIMESTAMP, INTERVAL $2 DAY);"
}


# $1 gnuplot script
# $2 is data file
# $3 is output file
function plotdata() {
	cat << EOF > $1
set datafile separator ","
set term svg enhanced mouse size 600,400 dynamic enhanced fname 'arial' fsize 10 name "Last_24H" butt 
set output '$3'

set key horizontal above

set xdata time
set timefmt '%Y-%m-%d %H:%M:%S'
set xlabel "Last Day's worth of data"
set format x "%H"
set autoscale x

set ylabel "pressure (mbar)"
set y2label "temperature %RH voltages"
set yrange [820:860]
set autoscale y2
set y2tics auto
set ytics auto

set grid ytics lt 0 lw 1 lc rgb "#880000"
set grid xtics lt 0 lw 1 lc rgb "#880000"

plot [:][:] "$2" using 1:5 title "pressure" with lines axis x1y1, "" using 1:2 title "Vdc" with lines axis x1y2, "" using 1:3 title "SHT11 RH" with lines axis x1y2, "" using 1:4 title "Internal RH" with lines axis x1y2, "" using 1:6 title "temperature #1" with lines axis x1y2, "" using 1:7 title "temperature #2" with lines axis x1y2

EOF
	gnuplot $1
}


extract $OUTFILE 1 $USER $PASSWORD $DATABASE
plotdata $GPLTFile $OUTFILE plot.svg