import { Component, AfterViewInit } from '@angular/core';
import * as L from 'leaflet'
import { HttpClient } from "@angular/common/http";
import { Observable, firstValueFrom } from 'rxjs';
import { Injectable } from  '@angular/core';
import { ɵafterNextNavigation } from '@angular/router';

export interface Telemetry {
  services: string[];
  longitude: number;
  latitude: number;
  ip_address: string;
  mainflux_version: string;
  last_seen: any;
  date: any;
  country: string;
  city: string;
}

export interface TelemetryPage {
  total: number;
  offset: number;
  limit: number;
  telemetry: Telemetry[];
}

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.css']
})
export class MapComponent implements AfterViewInit {
  private map!: L.Map;
  private limit: number = 10;
  private offset: number = 0;
  
  private initMap(): void {
    this.map = L.map('map').setView([20,0],2);
    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'}).addTo(this.map);
  }
  private markPins(telemetry: Telemetry[]) {
    telemetry.forEach(tel => {
      L.marker([tel.latitude, tel.longitude]).bindPopup(
        `<h3>Deployment details</h3>
        <p>IP Address:\t${tel.ip_address}</p>
        <p>Version:\t${tel.mainflux_version}</p>
        <p>Last seen:\t${tel.last_seen}</p>
        <p>Country:\t${tel.country}</p>
        <p>City:\t${tel.city}</p>
        <p>Services:\t${tel.services.join(',')}</p>`
      ).addTo(this.map);
    });
  }
  private async getPage() {
    while (true) {
      let tel = await this.telemetryService.retrieveTelemetry(this.limit,this.offset)
      this.offset = this.offset + tel.total
      if (tel != null) {
        if (tel.total==0) {
          break
        }
        this.markPins(tel.telemetry);
      }
    }
  }
  constructor(private telemetryService: TelemetryService) { }
  ngAfterViewInit(): void {
    this.initMap();
    this.getPage();
    
  }
}

@Injectable({
  providedIn:  'root'
})

export class TelemetryService {
  constructor(private httpClient: HttpClient) {};
  async retrieveTelemetry(limit: number, offset: number): Promise<TelemetryPage> {
    const getItems$: Observable<TelemetryPage> = this.httpClient.get<TelemetryPage>(`https://localhost/telemetry/sheets`, {
      headers: {'apikey':'77e04a7c-f207-40dd-8950-c344871fd516'},
      params: {"limit": limit, "offset": offset}
    }).pipe();
    return firstValueFrom(getItems$)
  }
}