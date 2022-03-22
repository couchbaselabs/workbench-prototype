import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, ObservableInput, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { AuthService } from './auth.service';
import { InitService } from './init.service';
import {
  Cluster,
  ClusterAddRequest,
  ClusterStatusResults,
  CloudClusterPage,
  CloudClusterHealth
} from './types';

@Injectable({
  providedIn: 'root',
})
export class RestService {
  constructor(
    private authService: AuthService,
    private initService: InitService,
    private http: HttpClient,
    private router: Router
  ) { }

  getClusters(): Observable<Cluster[]> {
    return this.http
      .get<Cluster[]>('/api/v1/clusters', { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  getCluster(uuid: string): Observable<Cluster> {
    return this.http
      .get<Cluster>(`/api/v1/clusters/${uuid}`, { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  getLogFile(uuid: string, node: string, logFile: string): Observable<any> {
    return this.http
      .get<any>(`/api/v1/clusters/${uuid}/nodes/${encodeURIComponent(node)}/logs/${logFile}`, {
        headers: this.getHeaders(),
        responseType: 'text' as 'json',
      })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  registerCluster(clusterAddReq: ClusterAddRequest): Observable<any> {
    return this.http
      .post('/api/v1/clusters', clusterAddReq, { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  deleteCluster(uuid: string): Observable<any> {
    return this.http
      .delete(`/api/v1/clusters/${uuid}`, { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  editCluster(uuid: string, request: any): Observable<any> {
    return this.http
      .patch(`/api/v1/clusters/${uuid}`, request, {
        headers: this.getHeaders(),
      })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  getCloudClusters(pagination: any): Observable<CloudClusterPage> {
    const httpParams = { params: new HttpParams({ fromObject: pagination }), headers: this.getHeaders() };
    return this.http.get<CloudClusterPage>('/api/v1/cloud/clusters', httpParams)
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  getCloudClusterHealth(id: string): Observable<CloudClusterHealth> {
    return this.http.get<CloudClusterHealth>(`/api/v1/cloud/clusters/${id}`, { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  addCloudCreds(creds: any): Observable<any> {
    return this.http.post('/api/v1/cloud/credentials', creds, { headers: this.getHeaders() })
      .pipe(catchError((err) => this.handleServerErrors(err)));
  }

  private handleServerErrors(err: HttpErrorResponse): ObservableInput<any> {
    switch (err.status) {
      case 401:
        this.authService.logOut();
        this.router.navigate(['login'], { queryParamsHandling: 'preserve' });
        break;
      case 503:
        this.initService.handle503();
        break;
    }

    return throwError(err);
  }

  private getHeaders(): HttpHeaders {
    const token = this.authService.getToken();
    return new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    });
  }
}
