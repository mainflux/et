import { Component } from '@angular/core';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css'],
})
export class SidebarComponent {
  countries: string[] = ["kenya", "serbia"]
  deployments: number = 0;
  noCountries: number = 0;
}
