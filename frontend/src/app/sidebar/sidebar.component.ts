import { Component, AfterViewInit } from '@angular/core';
import { Injectable } from  '@angular/core';
import { HttpClient } from "@angular/common/http";
import { Observable, firstValueFrom } from 'rxjs';

export interface TelemetrySummary {
  countries: string[];
  ip_addresses: string[];
}

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css'],
})

export class SidebarComponent implements AfterViewInit{
  countries!: string[];
  deployments!: number;
  noCountries!: number;
  private async getData() {
    let tel = await this.telemetryService.retrieveTelemetry()
    this.countries = tel.countries
    this.deployments = tel.ip_addresses.length;
    this.noCountries = tel.countries.length;
  }
  constructor(private telemetryService: TelemetryService) { }
  ngAfterViewInit(): void {
    this.getData()
  }
}

@Injectable({
  providedIn:  'root'
})

export class TelemetryService {
  constructor(private httpClient: HttpClient) {};
  async retrieveTelemetry(): Promise<TelemetrySummary> {
    const getItems$: Observable<TelemetrySummary> = this.httpClient.get<TelemetrySummary>(`https://localhost/telemetry/summary`, {
      headers: {'apikey':'77e04a7c-f207-40dd-8950-c344871fd516'},
    }).pipe();
    return firstValueFrom(getItems$)
  }
}
