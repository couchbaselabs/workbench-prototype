import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { MainComponent } from './main/main.component';
import { ClustersComponent } from './clusters/clusters.component';
import { NodesTableComponent } from './clusters/nodes-table/nodes-table.component';
import { MnElementCraneModule } from './mn.element.crane';
import { AddClusterPanelComponent } from './clusters/add-cluster-panel/add-cluster-panel.component';
import { BackButtonComponent } from './back-button.component';
import { EditClusterPanelComponent } from './clusters/edit-cluster-panel/edit-cluster-panel.component';
import { InitializeComponent } from './initialize/initialize.component';
import { FormatVersionPipe } from './format-version.pipe';
import { ClusterSubNavComponent } from './cluster-sub-nav/cluster-sub-nav.component';
import { OnpremComponent } from './clusters/onprem/onprem.component';
import { CloudComponent } from './clusters/cloud/cloud.component';
import { AddCredsComponent } from './clusters/cloud/add-creds/add-creds.component';
import { ClusterHealthComponent } from './clusters/cloud/cluster-health/cluster-health.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    MainComponent,
    ClustersComponent,
    NodesTableComponent,
    AddClusterPanelComponent,
    BackButtonComponent,
    EditClusterPanelComponent,
    InitializeComponent,
    FormatVersionPipe,
    ClusterSubNavComponent,
    OnpremComponent,
    CloudComponent,
    AddCredsComponent,
    ClusterHealthComponent
  ],
  imports: [
    BrowserModule,
    MnElementCraneModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    AppRoutingModule,
    NgbModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
