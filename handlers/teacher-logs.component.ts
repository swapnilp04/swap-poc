import { Component, OnInit } from '@angular/core';
import { TeacherLogService } from './../../service/teacher-log.service';
import { LoginService } from './../../service/login.service';
import { TeacherLog } from './../../interface/teacher-log';

@Component({
  selector: 'app-teacher-logs',
  templateUrl: './teacher-logs.component.html',
  styleUrls: ['./teacher-logs.component.css']
})

export class TeacherLogsComponent implements OnInit {
  
  public teacherLogs: TeacherLog[] = [];

  constructor(private teacherLogService: TeacherLogService, private loginService: LoginService){}

  ngOnInit(): void {
    this.loadTeacherLogs();
  }

  errorHandle(error: any): void {
    if(error.status == 401) {
      this.loginService.toLogin();
    }
    if(error.status == 0) {
      this.loginService.toLogin();
    }
  }

  loadTeacherLogs(): void {
    this.teacherLogService.getTeacherLogs().subscribe (
      (response: any) => this.assignTeacherLogs(response),
      (error: any) => this.errorHandle(error),
      () => console.log('Done getting TeacherLogs......')
    );
  }

  assignTeacherLogs(response: any) {
    this.teacherLogs = response
  }
}
