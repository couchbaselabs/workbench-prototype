import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './auth.guard';
import { LoginComponent } from './login/login.component';
import { MainComponent } from './main/main.component';
import { ClustersComponent } from './clusters/clusters.component'
import { InitializeComponent } from './initialize/initialize.component';
import { InitGuard } from './init.guard';

const routes: Routes = [
  { path: 'init', component: InitializeComponent, pathMatch: 'full'},
  { path: 'login', component: LoginComponent, pathMatch: 'full' , canActivate: [InitGuard]},
  { path: '', component: MainComponent, canActivate: [InitGuard, AuthGuard], children: [
    { path: 'clusters', component: ClustersComponent, canActivateChild: [InitGuard, AuthGuard]}
  ]},
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes),
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
