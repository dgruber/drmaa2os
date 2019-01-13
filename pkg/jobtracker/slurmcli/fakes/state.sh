#!/bin/sh
case $4 in 
  1.batch* ) echo RUNNING ;;
  2.batch* ) echo COMPLETING ;;
  3.batch* ) echo COMPLETED ;;
  4.batch* ) echo CANCELLED ;;
  5.batch* ) echo BOOT_FAIL ;;
  6.batch* ) echo CONFIGURING ;;
  7.batch* ) echo DEADLINE ;;
  8.batch* ) echo FAILED ;;
  9.batch* ) echo NODE_FAIL ;;
  10.batch* ) echo OUT_OF_MEMORY ;;
  11.batch* ) echo PENDING ;;
  12.batch* ) echo PREEMPTED ;;
  13.batch* ) echo RESV_DEL_HOLD ;;
  14.batch* ) echo REQUEUE_FED ;;
  15.batch* ) echo REQUEUE_HOLD ;;
  16.batch* ) echo REQUEUED ;;
  17.batch* ) echo RESIZING ;;
  18.batch* ) echo REVOKED ;;
  19.batch* ) echo SIGNALING ;;
  20.batch* ) echo SPECIAL_EXIT ;;
  21.batch* ) echo STAGE_OUT ;;
  22.batch* ) echo STOPPED ;;
  23.batch* ) echo SUSPENDED ;;
  24.batch* ) echo TIMEOUT ;;
  99.batch* ) echo WrOnGsTaTe;;
esac
echo ""

