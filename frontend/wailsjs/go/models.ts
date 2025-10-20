export namespace main {
	
	export class CreateDocumentWithPermissionsRequest {
	    naziv_dokumenta: string;
	    projekat_id: number;
	    radni_tok_id: number;
	    opis?: string;
	    rok?: string;
	    korisnici_dozvole: number[];
	
	    static createFrom(source: any = {}) {
	        return new CreateDocumentWithPermissionsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.naziv_dokumenta = source["naziv_dokumenta"];
	        this.projekat_id = source["projekat_id"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.opis = source["opis"];
	        this.rok = source["rok"];
	        this.korisnici_dozvole = source["korisnici_dozvole"];
	    }
	}
	export class CreateWorkflowWithPhasesRequest {
	    naziv: string;
	    tip_toka: string;
	    faze: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateWorkflowWithPhasesRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.naziv = source["naziv"];
	        this.tip_toka = source["tip_toka"];
	        this.faze = source["faze"];
	    }
	}

}

export namespace models {
	
	export class ActivityLogRequest {
	    tip_aktivnosti: string;
	    entitet_tip?: string;
	    entitet_id?: number;
	    naziv_entiteta?: string;
	    opis?: string;
	    rezultat?: string;
	    dodatne_informacije?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ActivityLogRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tip_aktivnosti = source["tip_aktivnosti"];
	        this.entitet_tip = source["entitet_tip"];
	        this.entitet_id = source["entitet_id"];
	        this.naziv_entiteta = source["naziv_entiteta"];
	        this.opis = source["opis"];
	        this.rezultat = source["rezultat"];
	        this.dodatne_informacije = source["dodatne_informacije"];
	    }
	}
	export class DocumentPermissionRequest {
	    dokument_id: number;
	    korisnik_id: number;
	    moze_citati: boolean;
	    moze_menjati: boolean;
	    moze_brisati: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DocumentPermissionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dokument_id = source["dokument_id"];
	        this.korisnik_id = source["korisnik_id"];
	        this.moze_citati = source["moze_citati"];
	        this.moze_menjati = source["moze_menjati"];
	        this.moze_brisati = source["moze_brisati"];
	    }
	}
	export class DocumentPermissionResponse {
	    dozvola_id: number;
	    dokument_id: number;
	    korisnik_id: number;
	    korisnicko_ime: string;
	    ime: string;
	    prezime: string;
	    moze_citati: boolean;
	    moze_menjati: boolean;
	    moze_brisati: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DocumentPermissionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dozvola_id = source["dozvola_id"];
	        this.dokument_id = source["dokument_id"];
	        this.korisnik_id = source["korisnik_id"];
	        this.korisnicko_ime = source["korisnicko_ime"];
	        this.ime = source["ime"];
	        this.prezime = source["prezime"];
	        this.moze_citati = source["moze_citati"];
	        this.moze_menjati = source["moze_menjati"];
	        this.moze_brisati = source["moze_brisati"];
	    }
	}
	export class Dokumenti {
	    dokument_id: number;
	    projekat_id?: number;
	    naziv_dokumenta: string;
	    folder_id?: number;
	    opis?: string;
	    tip_dokumenta?: string;
	    jezik_dokumenta?: string;
	    kljucne_reci?: string;
	    radni_tok_id?: number;
	    trenutna_faza_id?: number;
	    kreirao_korisnik_id: number;
	    // Go type: time
	    datuma_postavke: any;
	    // Go type: time
	    poslednja_izmena?: any;
	    naziv_projekta?: string;
	    ime_kreirao?: string;
	    naziv_faze?: string;
	    broj_verzija?: number;
	    progres?: number;
	
	    static createFrom(source: any = {}) {
	        return new Dokumenti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dokument_id = source["dokument_id"];
	        this.projekat_id = source["projekat_id"];
	        this.naziv_dokumenta = source["naziv_dokumenta"];
	        this.folder_id = source["folder_id"];
	        this.opis = source["opis"];
	        this.tip_dokumenta = source["tip_dokumenta"];
	        this.jezik_dokumenta = source["jezik_dokumenta"];
	        this.kljucne_reci = source["kljucne_reci"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.trenutna_faza_id = source["trenutna_faza_id"];
	        this.kreirao_korisnik_id = source["kreirao_korisnik_id"];
	        this.datuma_postavke = this.convertValues(source["datuma_postavke"], null);
	        this.poslednja_izmena = this.convertValues(source["poslednja_izmena"], null);
	        this.naziv_projekta = source["naziv_projekta"];
	        this.ime_kreirao = source["ime_kreirao"];
	        this.naziv_faze = source["naziv_faze"];
	        this.broj_verzija = source["broj_verzija"];
	        this.progres = source["progres"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Faze {
	    faza_id: number;
	    radni_tok_id: number;
	    naziv_faze: string;
	    redosled: number;
	
	    static createFrom(source: any = {}) {
	        return new Faze(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.faza_id = source["faza_id"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.naziv_faze = source["naziv_faze"];
	        this.redosled = source["redosled"];
	    }
	}
	export class IstorijaFazaDokumenta {
	    istorija_id: number;
	    dokument_id: number;
	    prethodna_faza_id?: number;
	    nova_faza_id: number;
	    korisnik_id: number;
	    // Go type: time
	    datum_promene: any;
	
	    static createFrom(source: any = {}) {
	        return new IstorijaFazaDokumenta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.istorija_id = source["istorija_id"];
	        this.dokument_id = source["dokument_id"];
	        this.prethodna_faza_id = source["prethodna_faza_id"];
	        this.nova_faza_id = source["nova_faza_id"];
	        this.korisnik_id = source["korisnik_id"];
	        this.datum_promene = this.convertValues(source["datum_promene"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Korisnici {
	    korisnik_id: number;
	    korisnicko_ime: string;
	    email: string;
	    ime?: string;
	    prezime?: string;
	    uloga_id: number;
	    status: string;
	    poslednja_prijava?: string;
	    kreiran_datuma: string;
	    naziv_uloge?: string;
	
	    static createFrom(source: any = {}) {
	        return new Korisnici(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.korisnik_id = source["korisnik_id"];
	        this.korisnicko_ime = source["korisnicko_ime"];
	        this.email = source["email"];
	        this.ime = source["ime"];
	        this.prezime = source["prezime"];
	        this.uloga_id = source["uloga_id"];
	        this.status = source["status"];
	        this.poslednja_prijava = source["poslednja_prijava"];
	        this.kreiran_datuma = source["kreiran_datuma"];
	        this.naziv_uloge = source["naziv_uloge"];
	    }
	}
	export class Projekti {
	    projekat_id: number;
	    naziv_projekta: string;
	    opis?: string;
	    datum_pocetka?: string;
	    datum_zavrsetka?: string;
	    status: string;
	    rukovodilac_id?: number;
	    radni_tok_id?: number;
	    rukovodilac_ime?: string;
	    broj_zadataka?: number;
	    broj_clanova?: number;
	
	    static createFrom(source: any = {}) {
	        return new Projekti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projekat_id = source["projekat_id"];
	        this.naziv_projekta = source["naziv_projekta"];
	        this.opis = source["opis"];
	        this.datum_pocetka = source["datum_pocetka"];
	        this.datum_zavrsetka = source["datum_zavrsetka"];
	        this.status = source["status"];
	        this.rukovodilac_id = source["rukovodilac_id"];
	        this.radni_tok_id = source["radni_tok_id"];
	        this.rukovodilac_ime = source["rukovodilac_ime"];
	        this.broj_zadataka = source["broj_zadataka"];
	        this.broj_clanova = source["broj_clanova"];
	    }
	}
	export class RadniTokovi {
	    radni_tok_id: number;
	    naziv: string;
	    tip_toka: string;
	    opis?: string;
	    da_li_je_sablon: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RadniTokovi(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.radni_tok_id = source["radni_tok_id"];
	        this.naziv = source["naziv"];
	        this.tip_toka = source["tip_toka"];
	        this.opis = source["opis"];
	        this.da_li_je_sablon = source["da_li_je_sablon"];
	    }
	}
	export class SkornjeAktivnosti {
	    log_id: number;
	    korisnik_id?: number;
	    korisnik_ime: string;
	    tip_aktivnosti: string;
	    entitet_tip?: string;
	    entitet_id?: number;
	    naziv_entiteta?: string;
	    opis?: string;
	    rezultat: string;
	    kreiran_datuma: string;
	
	    static createFrom(source: any = {}) {
	        return new SkornjeAktivnosti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.log_id = source["log_id"];
	        this.korisnik_id = source["korisnik_id"];
	        this.korisnik_ime = source["korisnik_ime"];
	        this.tip_aktivnosti = source["tip_aktivnosti"];
	        this.entitet_tip = source["entitet_tip"];
	        this.entitet_id = source["entitet_id"];
	        this.naziv_entiteta = source["naziv_entiteta"];
	        this.opis = source["opis"];
	        this.rezultat = source["rezultat"];
	        this.kreiran_datuma = source["kreiran_datuma"];
	    }
	}
	export class StatistikaAktivnosti {
	    tip_aktivnosti: string;
	    broj_aktivnosti: number;
	    broj_korisnika: number;
	    danas: number;
	    ove_nedelje: number;
	    ovog_meseca: number;
	    poslednja_aktivnost: string;
	
	    static createFrom(source: any = {}) {
	        return new StatistikaAktivnosti(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tip_aktivnosti = source["tip_aktivnosti"];
	        this.broj_aktivnosti = source["broj_aktivnosti"];
	        this.broj_korisnika = source["broj_korisnika"];
	        this.danas = source["danas"];
	        this.ove_nedelje = source["ove_nedelje"];
	        this.ovog_meseca = source["ovog_meseca"];
	        this.poslednja_aktivnost = source["poslednja_aktivnost"];
	    }
	}
	export class StatistikaDokumenata {
	    ukupno_dokumenata: number;
	    novih_dokumenata_mesecno: number;
	    broj_autora: number;
	    broj_projekata_sa_dokumentima: number;
	    prosecno_verzija_po_dokumentu: number;
	
	    static createFrom(source: any = {}) {
	        return new StatistikaDokumenata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ukupno_dokumenata = source["ukupno_dokumenata"];
	        this.novih_dokumenata_mesecno = source["novih_dokumenata_mesecno"];
	        this.broj_autora = source["broj_autora"];
	        this.broj_projekata_sa_dokumentima = source["broj_projekata_sa_dokumentima"];
	        this.prosecno_verzija_po_dokumentu = source["prosecno_verzija_po_dokumentu"];
	    }
	}
	export class Tagovi {
	    tag_id: number;
	    naziv_taga: string;
	
	    static createFrom(source: any = {}) {
	        return new Tagovi(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_id = source["tag_id"];
	        this.naziv_taga = source["naziv_taga"];
	    }
	}
	export class Uloge {
	    uloga_id: number;
	    naziv_uloge: string;
	
	    static createFrom(source: any = {}) {
	        return new Uloge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uloga_id = source["uloga_id"];
	        this.naziv_uloge = source["naziv_uloge"];
	    }
	}
	export class UploadDocumentRequest {
	    naziv_dokumenta: string;
	    projekat_id?: number;
	    folder_id?: number;
	    opis: string;
	    tip_dokumenta: string;
	    jezik_dokumenta: string;
	    tagovi: string[];
	    kljucne_reci: string;
	
	    static createFrom(source: any = {}) {
	        return new UploadDocumentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.naziv_dokumenta = source["naziv_dokumenta"];
	        this.projekat_id = source["projekat_id"];
	        this.folder_id = source["folder_id"];
	        this.opis = source["opis"];
	        this.tip_dokumenta = source["tip_dokumenta"];
	        this.jezik_dokumenta = source["jezik_dokumenta"];
	        this.tagovi = source["tagovi"];
	        this.kljucne_reci = source["kljucne_reci"];
	    }
	}
	export class VerzijeDokumenata {
	    verzija_id: number;
	    dokument_id: number;
	    verzija_oznaka?: string;
	    putanja_do_fajla: string;
	    velicina_fajla_mb?: number;
	    postavio_korisnik_id: number;
	    // Go type: time
	    datuma_postavke: any;
	
	    static createFrom(source: any = {}) {
	        return new VerzijeDokumenata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.verzija_id = source["verzija_id"];
	        this.dokument_id = source["dokument_id"];
	        this.verzija_oznaka = source["verzija_oznaka"];
	        this.putanja_do_fajla = source["putanja_do_fajla"];
	        this.velicina_fajla_mb = source["velicina_fajla_mb"];
	        this.postavio_korisnik_id = source["postavio_korisnik_id"];
	        this.datuma_postavke = this.convertValues(source["datuma_postavke"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Zadacic {
	    zadacic_id: number;
	    dokument_id: number;
	    opis: string;
	    izvrsen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Zadacic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.zadacic_id = source["zadacic_id"];
	        this.dokument_id = source["dokument_id"];
	        this.opis = source["opis"];
	        this.izvrsen = source["izvrsen"];
	    }
	}
	export class ZahteviPromeneFaze {
	    zahtev_id: number;
	    zadatak_id?: number;
	    dokument_id?: number;
	    podnosilac_zahteva_id: number;
	    zahtevana_faza_id: number;
	    status: string;
	    komentar?: string;
	    // Go type: time
	    datum_kreiranja: any;
	
	    static createFrom(source: any = {}) {
	        return new ZahteviPromeneFaze(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.zahtev_id = source["zahtev_id"];
	        this.zadatak_id = source["zadatak_id"];
	        this.dokument_id = source["dokument_id"];
	        this.podnosilac_zahteva_id = source["podnosilac_zahteva_id"];
	        this.zahtevana_faza_id = source["zahtevana_faza_id"];
	        this.status = source["status"];
	        this.komentar = source["komentar"];
	        this.datum_kreiranja = this.convertValues(source["datum_kreiranja"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace services {
	
	export class LoginResponse {
	    user?: models.Korisnici;
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = this.convertValues(source["user"], models.Korisnici);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

